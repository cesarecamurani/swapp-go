package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"swapp-go/cmd/internal/application/service"
	"swapp-go/cmd/internal/utils"
	"time"
)

type PasswordResetHandler struct {
	ResetService *service.PasswordResetService
	UserService  *service.UserService
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

func NewPasswordResetHandler(resetService *service.PasswordResetService, userService *service.UserService) *PasswordResetHandler {
	return &PasswordResetHandler{
		ResetService: resetService,
		UserService:  userService,
	}
}

func (handler *PasswordResetHandler) RequestReset(context *gin.Context) {
	var request PasswordResetRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := handler.UserService.GetUserByEmail(request.Email)
	if err != nil || user == nil {
		context.JSON(http.StatusOK, gin.H{"message": "User not found, no reset token was created."})
		return
	}

	token, err := handler.ResetService.GenerateAndSaveToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate reset token"})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Reset token generated",
		"token":   token,
	})
}

func (handler *PasswordResetHandler) ResetPassword(context *gin.Context) {
	var request ResetPasswordRequest

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resetToken, err := handler.ResetService.ValidateToken(request.Token)
	if err != nil || resetToken == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
		return
	}

	if resetToken.ExpiresAt.Before(time.Now()) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
		return
	}

	user, err := handler.UserService.GetUserByID(resetToken.UserID)
	if err != nil || user == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	hashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	_, err = handler.UserService.UpdateUser(user.ID, map[string]interface{}{"password": hashedPassword})
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	err = handler.ResetService.DeleteToken(request.Token)
	if err != nil {
		log.Printf("Warning: failed to delete password reset token: %v", err)
	}

	context.JSON(http.StatusOK, gin.H{"message": "Password reset successfully!"})
}
