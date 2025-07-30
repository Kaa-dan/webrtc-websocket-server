package helpers

import (
	"log"
	"os"

	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetail struct {
	Email     string
	FirstName string
	Last_name string
	Uid       string
	User_type string
	jwt.Claims
}

type TokenManager struct {
	userCollection *mongo.Collection
	secretKey      string
}

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
