package ds

import (
	"errors"

	"cloud.google.com/go/datastore"
)

// Cache interface
type Cache interface {
	Get(*datastore.Key, interface{}) error
	GetMulti([]*datastore.Key, interface{}) error
	Set(*datastore.Key, interface{}) error
	SetMulti([]*datastore.Key, interface{}) error
	Del(*datastore.Key) error
	DelMulti([]*datastore.Key) error
}

// Errors
var (
	ErrCacheNotFound = errors.New("ds: cache not found")
)
