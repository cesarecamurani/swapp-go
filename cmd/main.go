package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/config"
	"swapp-go/cmd/internal/validators"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	validators.Init()

	migrate()

	router := gin.Default()
	db := config.GetDB()
	config.SetupRoutes(router, db)

	router.Static("/uploads", "./uploads")

	err := router.Run(":9000")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func migrate() {
	if err := config.DB.AutoMigrate(&models.UserModel{}); err != nil {
		log.Fatalf("failed to migrate UserModel: %v", err)
	}

	if err := config.DB.AutoMigrate(&models.PasswordResetModel{}); err != nil {
		log.Fatalf("failed to migrate PasswordResetTokenModel: %v", err)
	}

	if err := config.DB.AutoMigrate(&models.ItemModel{}); err != nil {
		log.Fatalf("failed to migrate ItemModel: %v", err)
	}
}
