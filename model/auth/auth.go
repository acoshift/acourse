package auth

// SignUp command
type SignUp struct {
	Email    string
	Password string

	Result string // user id
}

// SendPasswordResetEmail command
type SendPasswordResetEmail struct {
	Email string
}
