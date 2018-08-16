package appctx

import (
	"context"

	"github.com/acoshift/flash"
)

const (
	keyUserID      = "user_id"
	keyOpenIDState = "openid_state"
)

// SetUserID sets user id to session
func SetUserID(ctx context.Context, userID string) {
	getSession(ctx).Set(keyUserID, userID)
}

// GetUserID gets user id from session
func GetUserID(ctx context.Context) string {
	return getSession(ctx).GetString(keyUserID)
}

// SetOpenIDState sets open id state to session
func SetOpenIDState(ctx context.Context, state string) {
	getSession(ctx).Set(keyOpenIDState, state)
}

// DelOpenIDState deletes open id state from session
func DelOpenIDState(ctx context.Context) {
	getSession(ctx).Del(keyOpenIDState)
}

// GetOpenIDState gets open id state from session
func GetOpenIDState(ctx context.Context) string {
	return getSession(ctx).GetString(keyOpenIDState)
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
