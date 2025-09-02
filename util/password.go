package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// generate a hash for the password
func HashPassword(password string) (string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error while hashing: %w", err)
	}
	return string(hashedPassword), nil
}

// function to validate password
func CheckPassword (password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}