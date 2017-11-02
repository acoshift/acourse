package app

import (
	"github.com/acoshift/session"
)

type sessionKey int

const (
	_ sessionKey = iota
	keyUserID
	keyOpenIDSessionID
)

// SetUserID sets user id to session
func SetUserID(sess *session.Session, userID string) {
	sess.Set(keyUserID, userID)
}

// GetUserID gets user id from session
func GetUserID(sess *session.Session) string {
	id, _ := sess.Get(keyUserID).(string)
	return id
}

// SetOpenIDSessionID sets open id session id to session
func SetOpenIDSessionID(sess *session.Session, openIDSessionID string) {
	sess.Set(keyOpenIDSessionID, openIDSessionID)
}

// GetOpenIDSessionID gets open id session id from session
func GetOpenIDSessionID(sess *session.Session) string {
	id, _ := sess.Get(keyOpenIDSessionID).(string)
	return id
}
