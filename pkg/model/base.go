package model

import (
	"strconv"

	"cloud.google.com/go/datastore"
)

// Base is the database base model
type Base struct {
	key *datastore.Key
	ID  string `datastore:"-"`
}

// KeySetter interface
type KeySetter interface {
	SetKey(*datastore.Key)
}

// KeyGetSetter interface
type KeyGetSetter interface {
	SetKey(*datastore.Key)
	Key() *datastore.Key
	NewKey(string)
}

// SetKey sets key to model
func (x *Base) SetKey(key *datastore.Key) {
	x.key = key
	if key == nil {
		x.ID = ""
	} else {
		x.ID = key.Name
		if x.ID == "" {
			x.ID = strconv.FormatInt(key.ID, 10)
		}
	}
}

// Key expose key from model
func (x *Base) Key() *datastore.Key {
	return x.key
}

// NewKey sets incomplete key to key
func (x *Base) NewKey(kind string) {
	x.key = datastore.IncompleteKey(kind, nil)
}
