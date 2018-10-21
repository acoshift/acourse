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

// SignInPassword command
type SignInPassword struct {
	Email    string
	Password string

	Result string // user id
}

// GenerateOpenIDURI command
type GenerateOpenIDURI struct {
	Provider string

	Result struct {
		RedirectURI string
		State       string
	}
}

// SignInOpenIDCallback command
type SignInOpenIDCallback struct {
	URI   string
	State string

	Result string // user id
}
