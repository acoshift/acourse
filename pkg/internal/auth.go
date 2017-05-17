package internal

import (
	"crypto/rand"

	identitytoolkit "google.golang.org/api/identitytoolkit/v3"
)

// SignInUser sign in user with email and password
func SignInUser(email, password string) (string, error) {
	resp, err := gitClient.VerifyPassword(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
		Email:    email,
		Password: password,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

func generateSessionID() string {
	b := make([]byte, 24)
	rand.Read(b)
	return string(b)
}

// SignInUserProvider sign in user with open id provider
func SignInUserProvider(provider string) (redirectURI string, sessionID string, err error) {
	sessID := generateSessionID()
	resp, err := gitClient.CreateAuthUri(&identitytoolkit.IdentitytoolkitRelyingpartyCreateAuthUriRequest{
		ProviderId:   provider,
		ContinueUri:  baseURL + "/openid/callback",
		AuthFlowType: "CODE_FLOW",
		SessionId:    sessID,
	}).Do()
	if err != nil {
		return "", "", err
	}
	return resp.AuthUri, sessID, nil
}

// SignInUserProviderCallback sign in user with open id provider callback
func SignInUserProviderCallback(callbackURI string, sessID string) (string, error) {
	resp, err := gitClient.VerifyAssertion(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyAssertionRequest{
		RequestUri: baseURL + callbackURI,
		SessionId:  sessID,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

// SignUpUser creates new user
func SignUpUser(email, password string) (string, error) {
	resp, err := gitClient.SignupNewUser(&identitytoolkit.IdentitytoolkitRelyingpartySignupNewUserRequest{
		Email:    email,
		Password: password,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.LocalId, nil
}

// GetVerifyEmailCode gets out-of-band confirmation code for verify email
func GetVerifyEmailCode(email string) (string, error) {
	resp, err := gitClient.GetOobConfirmationCode(&identitytoolkit.Relyingparty{
		Kind:        "identitytoolkit#relyingparty",
		RequestType: "VERIFY_EMAIL",
		Email:       email,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.OobCode, nil
}

// GetResetPasswordCode gets out-of-band confirmation code for reset password
func GetResetPasswordCode(email string) (string, error) {
	resp, err := gitClient.GetOobConfirmationCode(&identitytoolkit.Relyingparty{
		Kind:        "identitytoolkit#relyingparty",
		RequestType: "PASSWORD_RESET",
		Email:       email,
	}).Do()
	if err != nil {
		return "", err
	}
	return resp.OobCode, nil
}
