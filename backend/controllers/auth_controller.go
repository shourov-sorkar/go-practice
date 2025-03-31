package controllers

import (
	"context"
	"fmt"
	"go-react-mvc/backend/database"
	"go-react-mvc/backend/models"
	"go-react-mvc/backend/utils"
	"net/http"
	"os"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(c *gin.Context) {
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
						errorMessages["email"] = "Please enter a valid email address"
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
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to register user", map[string]string{"error": err.Error()})
		return
	}

	collection := database.GetCollection("go_database", "users")

	if err := utils.CheckDuplicate(c, utils.CheckDuplicateParams{
		Model:      user,
		Collection: collection,
		Field:      "username",
		Value:      user.Username,
	}); err != nil {
		return
	}

	if err := utils.CheckDuplicate(c, utils.CheckDuplicateParams{
		Model:      user,
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
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to register user", map[string]string{"error": err.Error()})
		return
	}

	utils.SendSuccessResponse(c, http.StatusCreated, "User registered successfully", gin.H{
		"id":        user.ID,
		"name":      user.Name,
		"username":  user.Username,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	})
}

func Login(c *gin.Context) {
	var loginUser models.LoginUser
	if err := c.ShouldBindJSON(&loginUser); err != nil {
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
			}
			utils.SendErrorResponse(c, http.StatusBadRequest, "Validation failed", errorMessages)
			return
		}
		utils.SendErrorResponse(c, http.StatusBadRequest, "Failed to login", map[string]string{"error": err.Error()})
		return
	}

	collection := database.GetCollection("go_database", "users")

	var foundUser models.User
	err := collection.FindOne(context.TODO(), bson.M{"username": loginUser.Username}).Decode(&foundUser)
	if err != nil {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Username doesn't exist", map[string]string{"error": fmt.Sprintf("%s doesn't exist", loginUser.Username)})
		return
	}

	if !utils.ComparePasswords(foundUser.Password, loginUser.Password) {
		utils.SendErrorResponse(c, http.StatusUnauthorized, "Password doesn't match", map[string]string{"error": fmt.Sprintf("%s doesn't match", loginUser.Password)})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", map[string]string{"error": err.Error()})
		return
	}
	utils.SendSuccessResponse(c, http.StatusOK, "Login successful", gin.H{
		"username": loginUser.Username,
		"token":    tokenString,
	})
}
