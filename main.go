package main

import (
	"log"

	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/Kaa-dan/webrtc-websocket-server.git/handlers"
	"github.com/Kaa-dan/webrtc-websocket-server.git/helpers"
	"github.com/Kaa-dan/webrtc-websocket-server.git/managers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Loading env variabes from .env file
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	// connect to DB
	database.ConnectDB()
}

func main() {
	// router setup
	router := gin.Default()
	//logger middleware setup
	router.Use(gin.Logger())
	// cors setup
	router.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"http://localhost:3000"},
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		},
	))

	//auth route setup
	tokenHelper := helpers.NewTokenHelper()
	authManager := managers.NewAuthManager(tokenHelper)
	authHandler := handlers.NewAuthHandler(authManager)
	//registering auth-routes
	authHandler.RegisterAuthApis(router)

	//starting server
	router.Run(":8000")
}
