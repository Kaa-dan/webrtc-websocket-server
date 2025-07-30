package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBInstance struct {
	Client *mongo.Client
	DB     *mongo.Database
}

var DB DBInstance

func ConnectDB() {
	mongoURI := os.Getenv("MONGODB_URL")
	if mongoURI == "" {
		log.Fatal("DB url is not set")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "Chat-app" // default DB name
	}
	// Context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	// Set client option
	clientOptions := options.Client().ApplyURI(mongoURI)

	//Connect to DB
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal("Failed to connect to DB :", err)
	}
	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDb :", err)
	}

	fmt.Println("Connected to MongoDB")
	// Set the global DB instance
	DB = DBInstance{
		Client: client,
		DB:     client.Database(dbName),
	}
}

// GetCollection returns a collection from the database
func GetCollection(collectionName string) *mongo.Collection {
	return DB.DB.Collection(collectionName)
}

// DiconnectDb closes the MongoDb connection
func DisconnectDB() {
	if DB.Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := DB.Client.Disconnect(ctx); err != nil {
			log.Fatal("Failed to disconnect from MongoDb", err)
		}
		fmt.Println("Disconnected from mongoDB")
	}
}
