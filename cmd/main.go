package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/config"
)

func main() {
	config.InitDB()
	migrate()

	server := gin.Default()

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

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
