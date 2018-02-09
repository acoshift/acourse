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

// setUserID sets user id to session
func setUserID(sess *session.Session, userID string) {
	sess.Set(keyUserID, userID)
}

// getUserID gets user id from session
func getUserID(sess *session.Session) string {
	id, _ := sess.Get(keyUserID).(string)
	return id
}

// setOpenIDSessionID sets open id session id to session
func setOpenIDSessionID(sess *session.Session, openIDSessionID string) {
	sess.Set(keyOpenIDSessionID, openIDSessionID)
}

// delOpenIDSessionID deletes open id session id from session
func delOpenIDSessionID(sess *session.Session) {
	sess.Del(keyOpenIDSessionID)
}

// getOpenIDSessionID gets open id session id from session
func getOpenIDSessionID(sess *session.Session) string {
	id, _ := sess.Get(keyOpenIDSessionID).(string)
	return id
}
