package redis

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/acoshift/ds"
	"github.com/garyburd/redigo/redis"
)

// Cache implement Cache interface
type Cache struct {
	Pool      *redis.Pool
	Prefix    string
	TTL       time.Duration
	ExtendTTL bool
	Skip      func(*datastore.Key) bool
}

func encode(v interface{}) ([]byte, error) {
	w := &bytes.Buffer{}
	err := gob.NewEncoder(w).Encode(v)
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func decode(b []byte, v interface{}) error {
	return gob.NewDecoder(bytes.NewReader(b)).Decode(v)
}

// Get gets data
func (cache *Cache) Get(key *datastore.Key, dst interface{}) error {
	if cache.Skip != nil && cache.Skip(key) {
		return ds.ErrCacheNotFound
	}

	db := cache.Pool.Get()
	defer db.Close()
	k := cache.Prefix + key.String()
	b, err := redis.Bytes(db.Do("GET", k))
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return ds.ErrCacheNotFound
	}
	if cache.ExtendTTL {
		db.Do("EXPIRE", k, int(cache.TTL/time.Second))
	}
	return decode(b, dst)
}

// GetMulti gets multi data
func (cache *Cache) GetMulti(keys []*datastore.Key, dst interface{}) error {
	db := cache.Pool.Get()
	defer db.Close()
	for _, key := range keys {
		db.Send("GET", cache.Prefix+key.String())
	}
	err := db.Flush()
	if err != nil {
		return err
	}
	for i := range keys {
		b, err := redis.Bytes(db.Receive())
		if err == redis.ErrNil {
			continue
		}
		if err != nil {
			return err
		}
		if len(b) > 0 {
			decode(b, reflect.Indirect(reflect.ValueOf(dst)).Index(i).Addr().Interface())
		}
	}
	if cache.ExtendTTL {
		ttl := int(cache.TTL / time.Second)
		for _, key := range keys {
			db.Send("EXPIRE", cache.Prefix+key.String(), ttl)
		}
		db.Flush()
	}
	return nil
}

// Set sets data
func (cache *Cache) Set(key *datastore.Key, src interface{}) error {
	if key == nil || src == nil {
		return nil
	}
	if cache.Skip != nil && cache.Skip(key) {
		return nil
	}

	db := cache.Pool.Get()
	defer db.Close()
	b, err := encode(src)
	if err != nil {
		return err
	}
	if cache.TTL > 0 {
		_, err = db.Do("SETEX", cache.Prefix+key.String(), int(cache.TTL/time.Second), b)
		return err
	}
	_, err = db.Do("SET", cache.Prefix+key.String(), b)
	return err
}

// SetMulti sets data
func (cache *Cache) SetMulti(keys []*datastore.Key, src interface{}) error {
	db := cache.Pool.Get()
	defer db.Close()
	ttl := int(cache.TTL / time.Second)
	db.Send("MULTI")
	for i, key := range keys {
		if key == nil {
			continue
		}
		if cache.Skip != nil && cache.Skip(key) {
			continue
		}
		s := reflect.Indirect(reflect.ValueOf(src)).Index(i)
		if s.IsNil() {
			continue
		}
		b, err := encode(s.Interface())
		if err != nil {
			return err
		}
		if cache.TTL > 0 {
			db.Send("SETEX", cache.Prefix+key.String(), ttl, b)
		} else {
			db.Send("SET", cache.Prefix+key.String(), b)
		}
	}
	_, err := db.Do("EXEC")
	return err
}

// Del dels data
func (cache *Cache) Del(key *datastore.Key) error {
	if key == nil {
		return nil
	}
	if cache.Skip != nil && cache.Skip(key) {
		return nil
	}

	db := cache.Pool.Get()
	defer db.Close()
	_, err := db.Do("DEL", cache.Prefix+key.String())
	return err
}

// DelMulti dels multi data
func (cache *Cache) DelMulti(keys []*datastore.Key) error {
	db := cache.Pool.Get()
	defer db.Close()
	db.Send("MULTI")
	for _, key := range keys {
		if key == nil {
			continue
		}
		if cache.Skip != nil && cache.Skip(key) {
			continue
		}
		db.Send("DEL", cache.Prefix+key.String())
	}
	_, err := db.Do("EXEC")
	return err
}
