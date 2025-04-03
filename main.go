package main

import (
	"go-react-poc/database"
	"go-react-poc/routes"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func initializeDatabase() (*mongo.Client, error) {
	client, err := database.ConnectDB()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully connected to MongoDB")
	return client, nil
}

func setupServer() *gin.Engine {
	router := routes.SetupGinRoutes()
	return router
}

func main() {
	// Initialize database connection
	client, err := initializeDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DisconnectDB(client)

	// Setup and run server
	router := setupServer()
	router.Run(":8080")
}
