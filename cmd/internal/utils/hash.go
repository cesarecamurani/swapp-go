package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	encryptedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(encryptedPwd), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
