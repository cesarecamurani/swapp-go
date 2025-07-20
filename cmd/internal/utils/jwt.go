package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func getJwtSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not set in environment")
	}
	return []byte(secret)
}

func GenerateToken(email string, userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"sub":   userId,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(getJwtSecret())
}
