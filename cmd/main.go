package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/config"
)

func main() {
	config.LoadEnv()
	config.InitDB()
	migrate()

	server := gin.Default()
	db := config.GetDB()

	config.SetupRoutes(server, db)

	err := server.Run(":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func migrate() {
	err := config.DB.AutoMigrate(&persistence.UserModel{})
	if err != nil {
		return
	}
}
