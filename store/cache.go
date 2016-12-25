package store

import (
	"sync"
	"time"
)

// Cache type
type Cache struct {
	mutex *sync.RWMutex
	data  map[string]interface{}
	ts    map[string]int64
	ttl   int64
}

// NewCache creates new cache group
func NewCache(ttl time.Duration) *Cache {
	return &Cache{
		mutex: &sync.RWMutex{},
		data:  map[string]interface{}{},
		ts:    map[string]int64{},
		ttl:   int64(ttl),
	}
}

// Get retrieves data from an index
func (c *Cache) Get(index string) interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	if c.ts[index] > 0 && c.ts[index]+c.ttl < time.Now().UnixNano() {
		delete(c.data, index)
		delete(c.ts, index)
	}
	return c.data[index]
}

// Set writes data to an index
func (c *Cache) Set(index string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[index] = value
	c.ts[index] = time.Now().UnixNano()
}

// Del deletes cache data
func (c *Cache) Del(index string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.data, index)
	delete(c.ts, index)
}

// Purge deletes all data
func (c *Cache) Purge() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = map[string]interface{}{}
	c.ts = map[string]int64{}
}
