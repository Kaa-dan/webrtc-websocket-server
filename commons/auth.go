package commons

import "strings"

// signup

type SignupInput struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password *string `json:"password" validate:"required,min=6"`
}

func NewSignupInput() *SignupInput {
	return &SignupInput{}
}

// HandleValidationError handles input validation errors with predefined error messages
func HandleValidationError(userData *SignupInput) error {
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
	if IsValidEmail(userData.Email) {
		return ErrInvalidEmail
	}

	return nil
}
