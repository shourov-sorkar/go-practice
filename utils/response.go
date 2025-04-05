package utils

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func SendErrorResponse(c *gin.Context, statusCode int, message string, details ...map[string]string) {
	var detailsMap map[string]string
	if len(details) > 0 {
		detailsMap = details[0]
	}
	response := ErrorResponse{
		Status:  statusCode,
		Message: message,
		Details: detailsMap,
	}
	c.JSON(statusCode, response)
}

func SendSuccessResponse(c *gin.Context, statusCode int, message string, data ...interface{}) {
	var dataValue interface{}
	if len(data) > 0 {
		dataValue = data[0]
	}
	response := SuccessResponse{
		Status:  statusCode,
		Message: message,
		Details: dataValue,
	}
	c.JSON(statusCode, response)
}
