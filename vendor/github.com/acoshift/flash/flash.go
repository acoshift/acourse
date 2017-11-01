package flash

import (
	"bytes"
	"encoding/gob"
	"net/url"
)

// Flash type
type Flash struct {
	v       url.Values
	changed bool
}

// New creates new flash
func New() *Flash {
	return &Flash{v: make(url.Values)}
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
		return nil, err
	}
	return buf.Bytes(), nil
}

// Set sets value to flash
func (f *Flash) Set(key, value string) {
	if !f.changed {
		f.changed = true
	}
	f.v.Set(key, value)
}

// Get gets value from flash
func (f *Flash) Get(key string) string {
	return f.v.Get(key)
}

// Values returns clone of flash's values
func (f *Flash) Values() url.Values {
	return cloneValues(f.v)
}

// Add adds value to flash
func (f *Flash) Add(key, value string) {
	if !f.changed {
		f.changed = true
	}
	f.v.Add(key, value)
}

// Del deletes key from flash
func (f *Flash) Del(key string) {
	if !f.Has(key) {
		return
	}
	if !f.changed {
		f.changed = true
	}
	f.v.Del(key)
}

// Has checks is flash has a given key
func (f *Flash) Has(key string) bool {
	return len(f.v[key]) > 0
}

// Clear deletes all data
func (f *Flash) Clear() {
	for k := range f.v {
		if !f.changed {
			f.changed = true
		}
		f.v.Del(k)
	}
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

func cloneValues(src url.Values) url.Values {
	n := make(url.Values, len(src))
	for k, vv := range src {
		n[k] = make([]string, len(vv))
		for kk, vvv := range vv {
			n[k][kk] = vvv
		}
	}
	return n
}
