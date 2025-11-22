package http

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-scaffolding/internal/user/ports"
)

// RegisterUserRoutes registers all user routes
func RegisterUserRoutes(router *gin.Engine, userService ports.UserService) {
	handler := NewUserHandler(userService)

	// User routes
	users := router.Group("/users")
	{
		users.POST("", handler.CreateUser)
		users.GET("", handler.ListUsers)
		users.GET("/email/:email", handler.GetUserByEmail) // Must be before /:id to avoid route conflict
		users.GET("/:id", handler.GetUser)
		users.PUT("/:id", handler.UpdateUser)
		users.DELETE("/:id", handler.DeleteUser)
	}
}
