package flash

import (
	"bytes"
	"encoding/gob"
	"net/url"
)

// Flash type
type Flash url.Values

// New creates new flash
func New() Flash {
	return make(Flash)
}

// Decode decodes flash data
func Decode(b []byte) (Flash, error) {
	f := New()
	if len(b) == 0 {
		return f, nil
	}

	err := gob.NewDecoder(bytes.NewReader(b)).Decode(&f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Encode encodes flash
func (f Flash) Encode() ([]byte, error) {
	// empty flash encode to empty bytes
	if len(f) == 0 {
		return []byte{}, nil
	}

	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(f)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Set sets value to flash
func (f Flash) Set(key, value string) {
	url.Values(f).Set(key, value)
}

// Get gets value from flash
func (f Flash) Get(key string) string {
	return url.Values(f).Get(key)
}

// Add adds value to flash
func (f Flash) Add(key, value string) {
	url.Values(f).Add(key, value)
}

// Del deletes key from flash
func (f Flash) Del(key string) {
	url.Values(f).Del(key)
}

// Has checks is flash has a given key
func (f Flash) Has(key string) bool {
	return len(url.Values(f)[key]) > 0
}

// Clear deletes all data
func (f Flash) Clear() {
	v := url.Values(f)
	for k := range f {
		v.Del(k)
	}
}

// Clone clones flash
func (f Flash) Clone() Flash {
	v := url.Values(f)
	n := make(url.Values, len(v))
	for k, vv := range v {
		n[k] = make([]string, len(vv))
		for kk, vvv := range vv {
			n[k][kk] = vvv
		}
	}
	return Flash(n)
}
