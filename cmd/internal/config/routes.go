package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/adapters/middleware"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/application/service"
)

func SetupRoutes(server *gin.Engine, db *gorm.DB) {
	userRepo := persistence.NewGormUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	server.POST("/users/register", userHandler.RegisterUser)
	server.POST("/users/login", userHandler.LoginUser)

	// Protected routes
	protected := server.Group("/")
	protected.Use(middleware.JwtAuthMiddleware)
	{
		protected.GET("/users/:id", userHandler.GetUserByID)
	}
}
