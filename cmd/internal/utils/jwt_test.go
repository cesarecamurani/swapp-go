package utils_test

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"os"
	"swapp-go/cmd/internal/utils"
	"testing"
	"time"
)

const (
	email  = "test@email.com"
	userID = "6e9648ee-fd0b-4267-adcb-0c03b0176277"
)

func init() {
	_ = os.Setenv("JWT_SECRET", "test_jwt_secret")
}

func signingKey(_ *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET")), nil
}

func parseTestToken(t *testing.T, tokenString string) jwt.MapClaims {
	t.Helper()

	token, err := jwt.Parse(tokenString, signingKey)
	if err != nil || !token.Valid {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Token claims are not of type jwt.MapClaims")
	}

	return claims
}

func TestGenerateToken_Valid(t *testing.T) {
	tokenString, err := utils.GenerateToken(email, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	claims := parseTestToken(t, tokenString)

	assert.Equal(t, email, claims["email"])
	assert.Equal(t, userID, claims["sub"])
}

func TestGenerateToken_Expired(t *testing.T) {
	claims := jwt.MapClaims{
		"email": email,
		"sub":   "expired_user_id",
		"exp":   time.Now().Add(-time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	assert.NoError(t, err)

	parsedToken, err := jwt.Parse(tokenString, signingKey)

	assert.Error(t, err)
	assert.False(t, parsedToken.Valid)
}
