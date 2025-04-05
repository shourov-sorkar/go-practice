package controllers

import (
	"context"
	"go-react-poc/database"
	"go-react-poc/models"
	"go-react-poc/utils"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, e := range validationErrors {
				if e.Tag() == "required" {
					errorMsg := utils.RequiredFieldValidation(utils.RequiredFieldParams{
						Field: e.Field(),
						Value: e.Value().(string),
					})
					if errorMsg != "" {
						errorMessages[e.Field()] = errorMsg
					}
					continue
				}
				switch e.Field() {
				case "email":
					if e.Tag() == "email" {
						errorMessages["email"] = "Invalid email format"
					}
				case "password":
					if e.Tag() == "min" {
						errorMessages["password"] = "Password must be at least 6 characters long"
					}
				}
			}
			utils.SendErrorResponse(c, http.StatusBadRequest, "Validation failed", errorMessages)
			return
		}
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to create user", map[string]string{"error": err.Error()})
		return
	}

	collection := database.GetCollection("go_database", "users")

	if err := utils.CheckDuplicate(c, utils.CheckDuplicateParams{
		Collection: collection,
		Field:      "username",
		Value:      user.Username,
	}); err != nil {
		return
	}

	if err := utils.CheckDuplicate(c, utils.CheckDuplicateParams{
		Collection: collection,
		Field:      "email",
		Value:      user.Email,
	}); err != nil {
		return
	}

	user.Password = utils.HashPassword(user.Password)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create user", map[string]string{"error": err.Error()})
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "User created successfully", gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var existingUser models.User
	filter := bson.M{"_id": objectID}
	err = database.GetCollection("go_database", "users").FindOne(context.TODO(), filter).Decode(&existingUser)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	updateData := make(map[string]interface{})

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid request body", map[string]string{"error": err.Error()})
		return
	}

	if len(updateData) == 0 {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Request body cannot be empty")
		return
	}

	validFields := make(map[string]bool)
	userType := reflect.TypeOf(models.User{})
	for i := range userType.NumField() {
		field := userType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" {
			fieldName := strings.Split(jsonTag, ",")[0]
			if fieldName != "id" && fieldName != "createdAt" && fieldName != "updatedAt" {
				validFields[fieldName] = true
			}
		}
	}
	duplicateFields := make(map[string]string)
	for field, value := range updateData {
		if !validFields[field] {
			utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid field in request body", map[string]string{"error": field + " is not a valid field"})
			return
		}
		if field != "password" {
			existingValue := reflect.ValueOf(existingUser).FieldByNameFunc(func(name string) bool {
				f, _ := reflect.TypeOf(existingUser).FieldByName(name)
				jsonTag := f.Tag.Get("json")
				fieldName := strings.Split(jsonTag, ",")[0]
				return fieldName == field
			})

			if existingValue.IsValid() && value == existingValue.Interface() {
				duplicateFields[field] = "We already have this value for " + field
			}
		}
	}

	if len(duplicateFields) > 0 {
		utils.SendErrorResponse(c, http.StatusConflict, "Some fields have same values", duplicateFields)
		return
	}

	if password, ok := updateData["password"].(string); ok {
		updateData["password"] = utils.HashPassword(password)
	}

	updateData["updatedAt"] = time.Now()
	update := bson.M{"$set": updateData}

	updateResult, err := database.GetCollection("go_database", "users").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update user")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "User updated successfully", updateResult)
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	deleteResult, err := database.GetCollection("go_database", "users").DeleteOne(context.TODO(), bson.M{"_id": objectId})
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	if deleteResult.DeletedCount == 0 {
		utils.SendErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	utils.SendSuccessResponse(c, http.StatusOK, "User deleted successfully", deleteResult)
}
