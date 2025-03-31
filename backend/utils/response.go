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
	Data    interface{} `json:"data,omitempty"`
}

func SendErrorResponse(c *gin.Context, statusCode int, message string, details map[string]string) {
	response := ErrorResponse{
		Status:  statusCode,
		Message: message,
		Details: details,
	}
	c.JSON(statusCode, response)
}

func SendSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Status:  statusCode,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}
