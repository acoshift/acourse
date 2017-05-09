package internal

// contextKey is the internal acourse's context key type
// use for store value in context
type contextKey int

// context key values
const (
	_       contextKey = iota
	keyUser            // user object
)
