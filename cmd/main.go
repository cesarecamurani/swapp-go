package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/application/service"
	"swapp-go/cmd/internal/config"
)

func main() {
	config.InitDB()
	migrate()

	server := gin.Default()

	userRepo := persistence.NewGormUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	server.POST("/users/register", userHandler.RegisterUser)

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
