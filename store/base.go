package store

import (
	"strconv"

	"cloud.google.com/go/datastore"
)

// Base is the database base model
type Base struct {
	key *datastore.Key
	ID  string `datastore:"-"`
}

// isBase interface
type isBase interface {
	setKey(key *datastore.Key)
}

func (m *Base) setKey(key *datastore.Key) {
	m.key = key
	if key == nil {
		m.ID = ""
	} else {
		m.ID = key.Name
		if m.ID == "" {
			m.ID = strconv.FormatInt(key.ID, 10)
		}
	}
}
