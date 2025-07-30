package main

import (
	"log"

	"github.com/Kaa-dan/webrtc-websocket-server.git/database"
	"github.com/joho/godotenv"
)

func main() {
	// Loading env variabes from .env file
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found")
	}

	// connect to DB
	database.ConnectDB()

}
