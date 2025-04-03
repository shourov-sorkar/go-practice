package controllers

import (
	"context"
	"go-react-poc/database"
	"go-react-poc/models"
	"go-react-poc/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsers(c *gin.Context) {
	params := utils.GetPaginationParams(c, 10)

	totalCount, err := database.GetCollection("go_database", "users").CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to count users", map[string]string{"error": err.Error()})
		return
	}
	params.Total = totalCount

	options := options.Find().SetLimit(int64(params.Limit)).SetSkip(int64(params.Skip))
	cursor, err := database.GetCollection("go_database", "users").Find(context.TODO(), bson.M{}, options)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users", map[string]string{"error": err.Error()})
		return
	}
	defer cursor.Close(context.TODO())

	var users []models.UserResponse
	if err := cursor.All(context.TODO(), &users); err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to decode users", map[string]string{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.GetPaginatedResponse(users, params))
}

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	insertResult, err := database.GetCollection("go_database", "users").InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": insertResult.InsertedID,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	updateResult, err := database.GetCollection("go_database", "users").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"updated": updateResult.ModifiedCount,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid user ID format", map[string]string{"error": err.Error()})
		return
	}

	deleteResult, err := database.GetCollection("go_database", "users").DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", map[string]string{"error": err.Error()})
		return
	}

	if deleteResult.DeletedCount == 0 {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found", map[string]string{"error": "User not found"})
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "User deleted successfully", gin.H{
		"deleted": deleteResult.DeletedCount,
	})
}
