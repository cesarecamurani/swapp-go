package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"swapp-go/cmd/internal/application/service"
	"swapp-go/cmd/internal/domain"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type RegisterUserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (userHandler *UserHandler) RegisterUser(context *gin.Context) {
	var request RegisterUserRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	user := &domain.User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}

	err := userHandler.userService.RegisterUser(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := &RegisterUserResponse{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}

	context.JSON(http.StatusCreated, response)
}
