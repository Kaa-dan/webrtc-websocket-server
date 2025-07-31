package commons

// signup

type SignupInput struct {
	Username string  `json:"username" validate:"required"`
	Email    string  `json:"email" validate:"required"`
	Password *string `json:"password" validate:"required,min=6"`
}

func NewSignupInput() *SignupInput {
	return &SignupInput{}
}
