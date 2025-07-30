package helpers

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetail struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.RegisteredClaims
}

type TokenHelper struct {
	userCollection *mongo.Collection
	secretKey      string
}

// creates new TokenHelper instance
func NewTokenHelper() *TokenHelper {

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is env is not set")
	}
	return &TokenHelper{
		userCollection: database.GetCollection("users"),
		secretKey:      secretKey,
	}
}

// GenerateAllTokens  : generates both access and refresh token
func (tm *TokenHelper) GeneratAllToken(email, firstName, lastName, userType, uid string) (string, string, error) {
	// Create claims for access token (expires in 24 hours)
	claims := &SignedDetail{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// Create claims for refresh token (expires in 168 hours = 7 days)

	refreshClaims := &SignedDetail{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 168)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Generate access token
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(tm.secretKey))

	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(tm.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return token, refreshToken, nil
}

// ValidateToken validate the JWT  token and returns claims

func (tm *TokenHelper) ValidateToken(signedToken string) (*SignedDetail, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetail{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(tm.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*SignedDetail)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	return claims, nil
}

// UpdateAllTokens updates both access and refresh tokens in the database
func (tm *TokenHelper) UpdateAllTokens(signedToken, signedRefreshToken, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateObj := bson.D{
		{Key: "token", Value: signedToken},
		{Key: "refresh_token", Value: signedRefreshToken},
		{Key: "updated_at", Value: time.Now()},
	}

	filter := bson.M{"user_id": userId}
	update := bson.D{{Key: "$set", Value: updateObj}}
	opts := options.Update().SetUpsert(true)

	_, err := tm.userCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("failed to update tokens: %w", err)
	}

	return nil
}
