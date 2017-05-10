package app

type sessionKey int

const (
	_ sessionKey = iota
	keyUserID
	keyOpenIDSessionID
)
