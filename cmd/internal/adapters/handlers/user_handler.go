package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
	"net/http"
	"swapp-go/cmd/internal/adapters/handlers/responses"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/domain"
)

type UserHandler struct {
	userService services.UserServiceInterface
}

func NewUserHandler(userServiceInterface services.UserServiceInterface) *UserHandler {
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

func (handler *UserHandler) RegisterUser(context *gin.Context) {
	var request RegisterUserRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		responses.BadRequest(context, "Invalid request", err)
		return
	}

	user := &domain.User{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
		Phone:    request.Phone,
		Address:  request.Address,
	}

	if err := handler.userService.RegisterUser(user); err != nil {
		responses.BadRequest(context, "Invalid request", err)
		return
	}

	respondWithUser(context, http.StatusCreated, "User created successfully!", user)
}

func (handler *UserHandler) Update(context *gin.Context) {
	userID := context.GetString("userID")

	var request UpdateUserRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		responses.BadRequest(context, "Invalid request", err)
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		responses.BadRequest(context, "Invalid user ID", err)
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
			responses.BadRequest(context, "Invalid phone number", phoneErr)
			return
		}
		formattedPhone := phonenumbers.Format(parsed, phonenumbers.E164)
		updateData["phone"] = formattedPhone
	}
	if request.Address != nil {
		updateData["address"] = *request.Address
	}
	if len(updateData) == 0 {
		responses.BadRequest(context, "No valid fields provided for update", nil)
		return
	}

	updatedUser, err := handler.userService.Update(parsedID, updateData)
	if err != nil {
		responses.InternalServerError(context, "Failed to update user", err)
		return
	}

	respondWithUser(context, http.StatusOK, "User updated successfully!", updatedUser)
}

func (handler *UserHandler) Delete(context *gin.Context) {
	userID := context.GetString("userID")

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		responses.BadRequest(context, "Invalid user ID", err)
		return
	}

	if err = handler.userService.Delete(parsedID); err != nil {
		responses.InternalServerError(context, "Failed to delete user", err)
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "User deleted successfully!"})
}

func (handler *UserHandler) FindByID(context *gin.Context) {
	id := context.Param("id")

	userID, err := uuid.Parse(id)
	if err != nil {
		responses.BadRequest(context, "Invalid user ID", err)
		return
	}

	user, err := handler.userService.FindByID(userID)
	if err != nil {
		responses.NotFound(context, "User not found", err)
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

func (handler *UserHandler) LoginUser(context *gin.Context) {
	var request LoginUserRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		responses.BadRequest(context, "Invalid request", err)
		return
	}

	token, user, err := handler.userService.Authenticate(request.Username, request.Password)
	if err != nil {
		responses.Unauthorized(context, "Unauthorized", err)
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
