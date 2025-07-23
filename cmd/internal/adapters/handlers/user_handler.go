package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"swapp-go/cmd/internal/application/service"
	"swapp-go/cmd/internal/domain"
)

type UserHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(userServiceInterface service.UserServiceInterface) *UserHandler {
	return &UserHandler{userServiceInterface}
}

type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterUserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type LoginUserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
}

func (userHandler *UserHandler) RegisterUser(context *gin.Context) {
	var request RegisterUserRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	user := &domain.User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
	}

	err := userHandler.userService.RegisterUser(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	response := &RegisterUserResponse{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}

	context.JSON(http.StatusCreated, response)
}

func (userHandler *UserHandler) GetUserByID(context *gin.Context) {
	id := context.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID", "details": err.Error()})
		return
	}

	user, err := userHandler.userService.GetUserByID(userID)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "User not found", "details": err.Error()})
		return
	}

	response := &RegisterUserResponse{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}

	context.JSON(http.StatusOK, response)
}

func (userHandler *UserHandler) LoginUser(context *gin.Context) {
	var request LoginUserRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "details": err.Error()})
		return
	}

	token, user, err := userHandler.userService.Authenticate(request.Username, request.Password)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized", "details": err.Error()})
		return
	}

	response := &LoginUserResponse{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	context.JSON(http.StatusOK, response)
}
