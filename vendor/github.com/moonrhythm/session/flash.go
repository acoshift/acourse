package session

import (
	"bytes"
	"encoding/gob"
)

type flashData map[string][]interface{}

// Flash type
type Flash struct {
	v       flashData
	changed bool
}

func (f *Flash) decode(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	return gob.NewDecoder(bytes.NewReader(b)).Decode(&f.v)
}

// Encode encodes flash
func (f *Flash) encode() ([]byte, error) {
	// empty flash encode to empty bytes
	if len(f.v) == 0 {
		return []byte{}, nil
	}

	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(f.v)
	if err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// Values returns slice of given key
func (f *Flash) Values(key string) []interface{} {
	if !f.Has(key) {
		return []interface{}{}
	}

	f.changed = true
	r := f.v[key]
	delete(f.v, key)
	return r
}

// Set sets value to flash
func (f *Flash) Set(key string, value interface{}) {
	if !f.changed {
		f.changed = true
	}
	if f.v == nil {
		f.v = make(flashData)
	}
	f.v[key] = []interface{}{value}
}

// Add adds value to flash
func (f *Flash) Add(key string, value interface{}) {
	if !f.changed {
		f.changed = true
	}
	if f.v == nil {
		f.v = make(flashData)
	}
	f.v[key] = append(f.v[key], value)
}

// Get gets value from flash
func (f *Flash) Get(key string) interface{} {
	if !f.Has(key) {
		return nil
	}

	f.changed = true
	r := f.v[key][0]
	delete(f.v, key)
	return r
}

// GetString gets string from flash
func (f *Flash) GetString(key string) string {
	r, _ := f.Get(key).(string)
	return r
}

// GetInt gets int from flash
func (f *Flash) GetInt(key string) int {
	r, _ := f.Get(key).(int)
	return r
}

// GetInt64 gets int64 from flash
func (f *Flash) GetInt64(key string) int64 {
	r, _ := f.Get(key).(int64)
	return r
}

// GetFloat32 gets float32 from flash
func (f *Flash) GetFloat32(key string) float32 {
	r, _ := f.Get(key).(float32)
	return r
}

// GetFloat64 gets float64 from flash
func (f *Flash) GetFloat64(key string) float64 {
	r, _ := f.Get(key).(float64)
	return r
}

// GetBool gets bool from flash
func (f *Flash) GetBool(key string) bool {
	r, _ := f.Get(key).(bool)
	return r
}

// Del deletes key from flash
func (f *Flash) Del(key string) {
	if f.Has(key) {
		f.changed = true
	}
	delete(f.v, key)
}

// Has checks is flash has a given key
func (f *Flash) Has(key string) bool {
	if f.v == nil {
		return false
	}
	return len(f.v[key]) > 0
}

// Clear deletes all data
func (f *Flash) Clear() {
	if f.Count() > 0 {
		f.changed = true
	}
	f.v = nil
}

// Count returns count of flash's keys
func (f *Flash) Count() int {
	return len(f.v)
}

// Clone clones flash
func (f *Flash) Clone() *Flash {
	return &Flash{v: cloneValues(f.v)}
}

// Changed returns true if value changed
func (f *Flash) Changed() bool {
	return f.changed
}

func cloneValues(src flashData) flashData {
	n := make(flashData, len(src))
	for k, vv := range src {
		n[k] = make([]interface{}, len(vv))
		for kk, vvv := range vv {
			n[k][kk] = vvv
		}
	}
	return n
}
