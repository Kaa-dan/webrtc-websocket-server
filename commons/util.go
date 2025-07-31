package commons

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func IsValidEmail(email string) bool {
	// Basic email validation - you might want to use a more robust solution
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string, bCost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
