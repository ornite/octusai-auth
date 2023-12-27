package utils

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofor-little/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	database *mongo.Database // Exported database instance
)

// InitDB initializes a connection to the MongoDB database.
func InitDB() (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().ApplyURI(env.Get("MONGO_URI", ""))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to ensure connectivity
	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database = client.Database(env.Get("DATABASE", "octusai"))

	// Log that the database connection is established
	log.Println("Connected to MongoDB!")

	return database, nil
}

// GetDatabase returns the exported database instance.
func GetDatabase() *mongo.Database {
	return database
}
