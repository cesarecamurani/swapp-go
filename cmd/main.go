package main

import (
	"github.com/gin-gonic/gin"
	"log"
	modelsPkg "swapp-go/cmd/internal/adapters/persistence/models"
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
	models := []interface{}{
		&modelsPkg.UserModel{},
		&modelsPkg.PasswordResetModel{},
		&modelsPkg.ItemModel{},
		&modelsPkg.SwappRequestModel{},
	}

	if err := config.DB.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}
}
