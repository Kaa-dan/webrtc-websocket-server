package commons

// signup

type SignupInput struct {
	Fullname string `json:"full_name"`
}

func NewSignupInput() *SignupInput {
	return &SignupInput{}
}
