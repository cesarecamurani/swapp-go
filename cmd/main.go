package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/config"
	"swapp-go/cmd/internal/validators"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	migrate()
	validators.Init()

	router := gin.Default()
	db := config.GetDB()

	config.SetupRoutes(router, db)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func migrate() {
	if err := config.DB.AutoMigrate(&persistence.UserModel{}); err != nil {
		log.Fatalf("failed to migrate UserModel: %v", err)
	}

	if err := config.DB.AutoMigrate(&persistence.PasswordResetModel{}); err != nil {
		log.Fatalf("failed to migrate PasswordResetTokenModel: %v", err)
	}
}
