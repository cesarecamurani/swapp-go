package utils_test

import (
	"github.com/stretchr/testify/assert"
	"swapp-go/cmd/internal/utils"
	"testing"
)

var password = "super_secure_password123"

func TestHashPassword(t *testing.T) {
	hashedPassword, err := utils.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)
}

func TestCheckPasswordHash(t *testing.T) {
	hashedPassword, _ := utils.HashPassword(password)

	isValid := utils.CheckPasswordHash(password, hashedPassword)
	assert.True(t, isValid, "Password should be valid")

	isInvalid := utils.CheckPasswordHash(password, "invalid_hash")
	assert.False(t, isInvalid, "Password should not be valid")
}

func TestHashPasswordIsDifferentEachTime(t *testing.T) {
	hash1, _ := utils.HashPassword(password)
	hash2, _ := utils.HashPassword(password)

	assert.NotEqual(t, hash1, hash2, "Each hash should be unique due to salting")
}
