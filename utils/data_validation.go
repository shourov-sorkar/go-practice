package utils

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CheckDuplicateParams struct {
	Model      interface{}
	Collection *mongo.Collection
	Field      string
	Value      string
}
type RequiredFieldParams struct {
	Field string
	Value string
}

func Capitalize(str string) string {
	return cases.Title(language.English).String(str)
}

func CheckDuplicate(c *gin.Context, params CheckDuplicateParams) error {
	var existingData interface{}
	err := params.Collection.FindOne(context.TODO(), bson.M{params.Field: params.Value}).Decode(&existingData)
	if err == nil {
		SendErrorResponse(c, http.StatusConflict, Capitalize(params.Field)+" already exists", map[string]string{
			params.Field: params.Value + " already exists",
		})
		return fmt.Errorf("%s already exists", params.Field)
	}
	return nil
}

func RequiredFieldValidation(params RequiredFieldParams) string {
	if params.Value == "" {
		return params.Field + " is required"
	}
	return ""
}
