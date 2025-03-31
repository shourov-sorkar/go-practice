package middleware

import (
	"fmt"
	"go-react-mvc/backend/utils"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Authorization token required", map[string]string{"error": "Authorization token required"})
			c.Abort()
			return
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		fmt.Printf("Received token: %v\n", tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return secretKey, nil
		})

		if err != nil {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token", map[string]string{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			utils.SendErrorResponse(c, http.StatusUnauthorized, "Invalid token", map[string]string{"error": "Token is not valid"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", claims["user_id"])
		}

		c.Next()
	}
}
