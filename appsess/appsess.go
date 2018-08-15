package appsess

import (
	"github.com/acoshift/session"
)

const (
	keyUserID          = "user_id"
	keyOpenIDSessionID = "openid_session"
)

// SetUserID sets user id to session
func SetUserID(sess *session.Session, userID string) {
	sess.Set(keyUserID, userID)
}

// GetUserID gets user id from session
func GetUserID(sess *session.Session) string {
	return sess.GetString(keyUserID)
}

// SetOpenIDSessionID sets open id session id to session
func SetOpenIDSessionID(sess *session.Session, openIDSessionID string) {
	sess.Set(keyOpenIDSessionID, openIDSessionID)
}

// DelOpenIDSessionID deletes open id session id from session
func DelOpenIDSessionID(sess *session.Session) {
	sess.Del(keyOpenIDSessionID)
}

// GetOpenIDSessionID gets open id session id from session
func GetOpenIDSessionID(sess *session.Session) string {
	return sess.GetString(keyOpenIDSessionID)
}
