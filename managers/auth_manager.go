package managers

import (
	"context"
	"errors"
	"time"

	"github.com/Kaa-dan/webrtc-websocket-server.git/commons"
	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/Kaa-dan/webrtc-websocket-server.git/helpers"
	"github.com/Kaa-dan/webrtc-websocket-server.git/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthManager struct {
	tokenHelper    *helpers.TokenHelper
	userCollection *mongo.Collection
}

func NewAuthManager(token_Helper *helpers.TokenHelper) *AuthManager {
	return &AuthManager{
		tokenHelper:    token_Helper,
		userCollection: database.GetCollection("users"),
	}
}

func (aM *AuthManager) SignUp(user *commons.SignupInput) (*models.User, error) {
	// Ensure initialization
	if aM.userCollection == nil {
		return nil, errors.New("database connection not initialized")
	}

	// Validate input
	if user == nil {
		return nil, commons.ErrInvalidInput
	}

	// Additional validation using the commons function
	if err := commons.HandleValidationError(user); err != nil {
		return nil, err
	}

	// Check if user already exists by email
	var existingUser models.User
	err := aM.userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	} else if err != mongo.ErrNoDocuments {
		return nil, errors.New("database error while checking existing user")
	}

	// Check if username already exists
	err = aM.userCollection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return nil, errors.New("username already taken")
	} else if err != mongo.ErrNoDocuments {
		return nil, errors.New("database error while checking existing username")
	}

	// Hash the password
	if user.Password == nil || *user.Password == "" {
		return nil, errors.New("password is required")
	}

	hashedPassword, err := commons.HashPassword(*user.Password, 10)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create new user
	hashedPasswordStr := string(hashedPassword)
	userType := "USER" // Default user type
	userId := primitive.NewObjectID().Hex()

	newUser := &models.User{
		ID:            primitive.NewObjectID(),
		Username:      &user.Username,
		Email:         &user.Email,
		Password:      &hashedPasswordStr,
		Token:         nil, // Will be set during login
		User_type:     &userType,
		Refresh_token: nil, // Will be set during login
		Created_at:    time.Now(),
		Updated_at:    time.Now(),
		User_id:       userId,
	}

	// Insert user into database
	_, err = aM.userCollection.InsertOne(context.Background(), newUser)
	if err != nil {
		return nil, errors.New("failed to create user in database")
	}

	// Don't return password in response
	newUser.Password = nil

	return newUser, nil
}
