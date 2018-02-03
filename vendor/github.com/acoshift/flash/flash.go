package flash

import (
	"bytes"
	"encoding/gob"
)

// Flash type
type Flash struct {
	v       Data
	changed bool
}

// New creates new flash
func New() *Flash {
	return &Flash{}
}

// Decode decodes flash data
func Decode(b []byte) (*Flash, error) {
	f := New()
	if len(b) == 0 {
		return f, nil
	}

	err := gob.NewDecoder(bytes.NewReader(b)).Decode(&f.v)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Encode encodes flash
func (f *Flash) Encode() ([]byte, error) {
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

// Values returns flash's data
func (f *Flash) Values() Data {
	return f.v
}

// Value returns slice of given key
func (f *Flash) Value(key string) []interface{} {
	if f.v == nil {
		return []interface{}{}
	}
	if f.v[key] == nil {
		return []interface{}{}
	}
	return f.v[key]
}

// Set sets value to flash
func (f *Flash) Set(key string, value interface{}) {
	if !f.changed {
		f.changed = true
	}
	if f.v == nil {
		f.v = make(Data)
	}
	f.v.Set(key, value)
}

// Add adds value to flash
func (f *Flash) Add(key string, value interface{}) {
	if !f.changed {
		f.changed = true
	}
	if f.v == nil {
		f.v = make(Data)
	}
	f.v.Add(key, value)
}

// Get gets value from flash
func (f *Flash) Get(key string) interface{} {
	return f.v.Get(key)
}

// GetString gets string from flash
func (f *Flash) GetString(key string) string {
	return f.v.GetString(key)
}

// GetInt gets int from flash
func (f *Flash) GetInt(key string) int {
	return f.v.GetInt(key)
}

// GetInt64 gets int64 from flash
func (f *Flash) GetInt64(key string) int64 {
	return f.v.GetInt64(key)
}

// GetFloat32 gets float32 from flash
func (f *Flash) GetFloat32(key string) float32 {
	return f.v.GetFloat32(key)
}

// GetFloat64 gets float64 from flash
func (f *Flash) GetFloat64(key string) float64 {
	return f.v.GetFloat64(key)
}

// GetBool gets bool from flash
func (f *Flash) GetBool(key string) bool {
	return f.v.GetBool(key)
}

// Del deletes key from flash
func (f *Flash) Del(key string) {
	if f.Has(key) {
		f.changed = true
	}
	f.v.Del(key)
}

// Has checks is flash has a given key
func (f *Flash) Has(key string) bool {
	return f.v.Has(key)
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

func cloneValues(src Data) Data {
	n := make(Data, len(src))
	for k, vv := range src {
		n[k] = make([]interface{}, len(vv))
		for kk, vvv := range vv {
			n[k][kk] = vvv
		}
	}
	return n
}
