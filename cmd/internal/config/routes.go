package config

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/adapters/middleware"
	gormRepo "swapp-go/cmd/internal/adapters/persistence/gorm"
	"swapp-go/cmd/internal/application/services"
)

func SetupRoutes(server *gin.Engine, db *gorm.DB) {
	userRepo := gormRepo.NewUserGormRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	passwordResetRepo := gormRepo.NewPasswordResetGormRepository(db)
	passwordResetService := services.NewPasswordResetService(passwordResetRepo)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService, userService)

	itemRepo := gormRepo.NewItemGormRepository(db)
	itemService := services.NewItemService(itemRepo)
	itemHandler := handlers.NewItemHandler(itemService)

	server.POST("/users/register", userHandler.RegisterUser)
	server.POST("/users/login", userHandler.LoginUser)

	server.POST("/password-reset/request", passwordResetHandler.RequestReset)
	server.POST("/password-reset/reset", passwordResetHandler.ResetPassword)

	server.GET("/items/:id", itemHandler.GetItemByID)

	// Protected routes
	protected := server.Group("/")
	protected.Use(middleware.JwtAuthMiddleware(os.Getenv("JWT_SECRET")))
	{
		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PATCH("/users/update", userHandler.UpdateUser)
		protected.DELETE("/users/delete", userHandler.DeleteUser)
		protected.POST("/items/create", itemHandler.CreateItem)
		protected.PUT("/items/update/:id", itemHandler.UpdateItem)
		protected.DELETE("/items/delete/:id", itemHandler.DeleteItem)
	}
}
