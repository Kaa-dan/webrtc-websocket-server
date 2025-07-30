package helpers

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetail struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.RegisteredClaims
}

type TokenManager struct {
	userCollection *mongo.Collection
	secretKey      string
}

// creates new TokenManager instance
func NewTokenManager() *TokenManager {

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is env is not set")
	}
	return &TokenManager{
		userCollection: database.GetCollection("users"),
		secretKey:      secretKey,
	}
}

// GenerateAllTokens  : generates both access and refresh token
func (tm *TokenManager) GeneratAllToken(email, firstName, lastName, userType, uid string) (string, string, error) {
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
	token, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(tm.secretKey))

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

func (tm *TokenManager) ValidateToken(signedToken string) (*SignedDetail, error) {
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
