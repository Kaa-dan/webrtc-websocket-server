package managers

import (
	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/Kaa-dan/webrtc-websocket-server.git/helpers"
	"github.com/Kaa-dan/webrtc-websocket-server.git/models"
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

func (aM *AuthManager) SignUp(user *models.User) {

}
