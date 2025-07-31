package commons

import "errors"

// Input validation errors
var (
	ErrInvalidInput    = errors.New("invalid input data")
	ErrMissingUsername = errors.New("user name is required")
	ErrMissingEmail    = errors.New("email is required")
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrInvalidUserID   = errors.New("invalid user ID")
)
