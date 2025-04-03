package database

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Printf("No .env file found: %v\n", err)
	}
}

var Client *mongo.Client

// ConnectDB establishes a connection to MongoDB
func ConnectDB() (*mongo.Client, error) {
	// MongoDB connection URI
	uri := os.Getenv("MONGODB_URI")

	// Set client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully connected to MongoDB!")
	Client = client
	return client, nil
}

// DisconnectDB closes the MongoDB connection
func DisconnectDB(client *mongo.Client) error {
	err := client.Disconnect(context.TODO())
	if err != nil {
		return err
	}
	fmt.Println("Connection to MongoDB closed.")
	return nil
}

// GetCollection returns a handle to the specified collection
func GetCollection(dbName, collName string) *mongo.Collection {
	return Client.Database(dbName).Collection(collName)
}
