package commons

// signup

type SignupInput struct {
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password *string `json:"password" validate:"required,min=6"`
}

func NewSignupInput() *SignupInput {
	return &SignupInput{}
}
