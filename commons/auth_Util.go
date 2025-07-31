package commons

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using bcrypt
func HashPassword(password string, bCost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// HandleValidationError handles input validation errors with predefined error messages
func HandleValidationError(userData *SignupInput) error {
	fmt.Println("user_data", userData)
	if userData == nil {
		return ErrInvalidInput
	}

	if strings.TrimSpace(userData.Username) == "" {
		return ErrMissingUsername
	}

	if strings.TrimSpace(userData.Email) == "" {
		return ErrMissingEmail
	}

	// Add email format validation if needed
	if !IsValidEmail(userData.Email) {
		return ErrInvalidEmail
	}

	return nil
}

func IsValidEmail(email string) bool {
	// Basic email validation - you might want to use a more robust solution
	email = strings.TrimSpace(email)
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
