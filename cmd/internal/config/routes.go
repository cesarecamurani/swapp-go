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

	swapRequestRepo := gormRepo.NewSwapRequestGormRepository(db)
	swapRequestService := services.NewSwapRequestService(swapRequestRepo)
	swapRequestHandler := handlers.NewSwapRequestHandler(swapRequestService)

	// Public routes
	server.POST("/users/register", userHandler.RegisterUser)
	server.POST("/users/login", userHandler.LoginUser)
	server.POST("/password-reset/request", passwordResetHandler.RequestReset)
	server.POST("/password-reset/reset", passwordResetHandler.ResetPassword)
	server.GET("/items/:id", itemHandler.FindByID)

	// Protected routes
	protected := server.Group("/")
	protected.Use(middleware.JwtAuthMiddleware(os.Getenv("JWT_SECRET")))
	usersGroup := protected.Group("/users")
	{
		usersGroup.GET("/:id", userHandler.FindByID)
		usersGroup.PATCH("/update", userHandler.Update)
		usersGroup.DELETE("/delete", userHandler.Delete)
	}
	itemsGroup := protected.Group("/items")
	{
		itemsGroup.POST("/create", itemHandler.Create)
		itemsGroup.PUT("/update/:id", itemHandler.Update)
		itemsGroup.DELETE("/delete/:id", itemHandler.Delete)
	}
	swapRequestsGroup := protected.Group("/swap-requests")
	{
		swapRequestsGroup.POST("/create", swapRequestHandler.Create)
		swapRequestsGroup.GET("/:id", swapRequestHandler.FindByID)
		swapRequestsGroup.GET("/reference/:reference", swapRequestHandler.FindByReferenceNumber)
		swapRequestsGroup.GET("/list-by-user/:id", swapRequestHandler.ListByUser)
		swapRequestsGroup.GET("/list-by-status/:status", swapRequestHandler.ListByStatus)
		swapRequestsGroup.DELETE("/delete/:id", swapRequestHandler.Delete)
		swapRequestsGroup.PATCH("/update-status/:id", swapRequestHandler.UpdateStatus)
	}

}
