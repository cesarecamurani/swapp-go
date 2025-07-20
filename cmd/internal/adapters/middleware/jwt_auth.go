package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func JwtAuthMiddleware(context *gin.Context) {
	authHeader := context.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or malformed"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	claims := token.Claims.(jwt.MapClaims)
	context.Set("userId", claims["userId"])
	context.Set("email", claims["email"])

	context.Next()
}
