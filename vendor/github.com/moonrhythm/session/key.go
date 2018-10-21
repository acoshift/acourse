package session

const (
	// manager internal data
	timestampKey = "_session/timestamp"
	destroyedKey = "_session/destroyed" // for detect session hijack

	// session internal data
	flashKey = "_session/flash"
)
