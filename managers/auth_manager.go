package managers

import (
	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/Kaa-dan/webrtc-websocket-server.git/helpers"
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
