package appctx

import (
	"context"

	"github.com/acoshift/flash"
)

const (
	keyUserID          = "user_id"
	keyOpenIDSessionID = "openid_session"
)

// SetUserID sets user id to session
func SetUserID(ctx context.Context, userID string) {
	getSession(ctx).Set(keyUserID, userID)
}

// GetUserID gets user id from session
func GetUserID(ctx context.Context) string {
	return getSession(ctx).GetString(keyUserID)
}

// SetOpenIDSessionID sets open id session id to session
func SetOpenIDSessionID(ctx context.Context, openIDSessionID string) {
	getSession(ctx).Set(keyOpenIDSessionID, openIDSessionID)
}

// DelOpenIDSessionID deletes open id session id from session
func DelOpenIDSessionID(ctx context.Context) {
	getSession(ctx).Del(keyOpenIDSessionID)
}

// GetOpenIDSessionID gets open id session id from session
func GetOpenIDSessionID(ctx context.Context) string {
	return getSession(ctx).GetString(keyOpenIDSessionID)
}

// GetFlash gets flash from context
func GetFlash(ctx context.Context) *flash.Flash {
	return getSession(ctx).Flash()
}

// RegenerateSessionID regerates session id
func RegenerateSessionID(ctx context.Context) {
	getSession(ctx).Regenerate()
}

// DestroySession destroys session
func DestroySession(ctx context.Context) {
	getSession(ctx).Destroy()
}
