package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/adapters/infrastructure/email"
	"swapp-go/cmd/internal/adapters/middleware"
	gormRepo "swapp-go/cmd/internal/adapters/persistence/gorm"
	modelsPkg "swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/config"
	"swapp-go/cmd/internal/config/routes"
	"swapp-go/cmd/internal/validators"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	validators.Init()
	migrate()

	router := gin.Default()
	db := config.GetDB()

	userRepo := gormRepo.NewUserGormRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	passwordResetRepo := gormRepo.NewPasswordResetGormRepository(db)
	passwordResetService := services.NewPasswordResetService(passwordResetRepo)
	passwordResetHandler := handlers.NewPasswordResetHandler(passwordResetService, userService)

	itemRepo := gormRepo.NewItemGormRepository(db)
	itemService := services.NewItemService(itemRepo)
	itemHandler := handlers.NewItemHandler(itemService)

	emailConfig := config.LoadEmailConfig()
	var emailService = email.NewSmtpEmailService(emailConfig)

	swapRequestRepo := gormRepo.NewSwapRequestGormRepository(db)
	swapRequestService := services.NewSwapRequestService(swapRequestRepo, userRepo, itemRepo, emailService)
	swapRequestHandler := handlers.NewSwapRequestHandler(swapRequestService)

	routes.SetupRoutes(
		router,
		userHandler,
		itemHandler,
		swapRequestHandler,
		passwordResetHandler,
		middleware.JwtAuthMiddleware(os.Getenv("JWT_SECRET")),
	)

	router.Static("/uploads", "./uploads")

	err := router.Run(":9000")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func migrate() {
	models := []interface{}{
		&modelsPkg.UserModel{},
		&modelsPkg.PasswordResetModel{},
		&modelsPkg.ItemModel{},
		&modelsPkg.SwapRequestModel{},
	}

	if err := config.DB.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}
}
