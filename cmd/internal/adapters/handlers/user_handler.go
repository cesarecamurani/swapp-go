package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
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
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUserRequest struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
	Phone    *string `json:"phone,omitempty" binding:"omitempty,phone"`
	Address  *string `json:"address,omitempty"`
}

type UserResponse struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
}

type UserSuccessResponse struct {
	Message string        `json:"message"`
	User    *UserResponse `json:"user"`
}

type LoginUserResponse struct {
	UserID   string  `json:"user_id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Address  *string `json:"address"`
	Token    string  `json:"token"`
}

func (userHandler *UserHandler) RegisterUser(context *gin.Context) {
	var request RegisterUserRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		badRequestResponse(context, "Invalid request", err)
		return
	}

	user := &domain.User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
		Phone:    request.Phone,
		Address:  request.Address,
	}

	if err := userHandler.userService.RegisterUser(user); err != nil {
		badRequestResponse(context, "Invalid request", err)
		return
	}

	respondWithUser(context, http.StatusCreated, "User created successfully!", user)
}

func (userHandler *UserHandler) UpdateUser(context *gin.Context) {
	userID := context.GetString("userID")

	var request UpdateUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		badRequestResponse(context, "Invalid request", err)
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		badRequestResponse(context, "Invalid user ID", err)
		return
	}

	updateData := make(map[string]interface{})
	if request.Username != nil {
		updateData["username"] = *request.Username
	}
	if request.Email != nil {
		updateData["email"] = *request.Email
	}
	if request.Phone != nil {
		parsed, phoneErr := phonenumbers.Parse(*request.Phone, "")
		if phoneErr != nil || !phonenumbers.IsValidNumber(parsed) {
			badRequestResponse(context, "Invalid phone number", phoneErr)
			return
		}
		formattedPhone := phonenumbers.Format(parsed, phonenumbers.E164)
		updateData["phone"] = formattedPhone
	}
	if request.Address != nil {
		updateData["address"] = *request.Address
	}
	if len(updateData) == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No valid fields provided for update"})
		return
	}

	updatedUser, err := userHandler.userService.UpdateUser(parsedID, updateData)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	respondWithUser(context, http.StatusOK, "User updated successfully!", updatedUser)
}

func (userHandler *UserHandler) DeleteUser(context *gin.Context) {
	userID := context.GetString("userID")

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		badRequestResponse(context, "Invalid user ID", err)
		return
	}

	if err = userHandler.userService.DeleteUser(parsedID); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete user",
			"details": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully!"})
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

	response := &UserResponse{
		UserID:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Address:  user.Address,
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
		Phone:    user.Phone,
		Address:  user.Address,
		Token:    token,
	}

	context.JSON(http.StatusOK, response)
}

func respondWithUser(context *gin.Context, status int, message string, user *domain.User) {
	response := UserSuccessResponse{
		Message: message,
		User: &UserResponse{
			UserID:   user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Address:  user.Address,
		},
	}
	context.JSON(status, response)
}

func badRequestResponse(context *gin.Context, message string, err error) {
	context.JSON(http.StatusBadRequest, gin.H{
		"error":   message,
		"details": err.Error(),
	})
}
