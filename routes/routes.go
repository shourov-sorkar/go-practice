package routes

import (
	"go-react-poc/controllers"
	"go-react-poc/middleware"

	"github.com/gin-gonic/gin"
)

func SetupGinRoutes() *gin.Engine {
	router := gin.Default()

	// Auth routes (No auth required)
	router.POST("/api/auth/register", controllers.Register)
	router.POST("/api/auth/login", controllers.Login)

	// Protected routes (Require valid JWT)
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// User routes (Require authentication)
	protected.GET("/users", controllers.GetUsers)
	protected.POST("/users", controllers.CreateUser)
	protected.PUT("/users/:id", controllers.UpdateUser)
	protected.DELETE("/users/:id", controllers.DeleteUser)

	return router
}
