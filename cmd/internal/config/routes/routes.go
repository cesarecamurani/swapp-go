package routes

import (
	"github.com/gin-gonic/gin"
	"swapp-go/cmd/internal/adapters/handlers"
)

func SetupRoutes(
	server *gin.Engine,
	userHandler *handlers.UserHandler,
	itemHandler *handlers.ItemHandler,
	swapRequestHandler *handlers.SwapRequestHandler,
	passwordResetHandler *handlers.PasswordResetHandler,
	authMiddleware gin.HandlerFunc,
) {

	// Public routes
	server.POST("/users/register", userHandler.RegisterUser)
	server.POST("/users/login", userHandler.LoginUser)
	server.POST("/password-reset/request", passwordResetHandler.RequestReset)
	server.POST("/password-reset/reset", passwordResetHandler.ResetPassword)
	server.GET("/items/:id", itemHandler.FindByID)

	// Protected routes
	protected := server.Group("/")
	protected.Use(authMiddleware)
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
