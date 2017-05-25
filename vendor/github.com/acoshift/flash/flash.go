package flash

import (
	"bytes"
	"context"
	"encoding/gob"
	"net/http"
	"net/url"

	"github.com/acoshift/middleware"
	"github.com/acoshift/session"
)

// Flash type
type Flash url.Values

func init() {
	gob.Register(flashKey)
}

// Middleware decodes flash data from session and save back
func Middleware() middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			sess := session.Get(ctx)
			var f Flash
			if b, ok := sess.Get(flashKey).([]byte); ok {
				f, _ = Decode(b)
			}
			if f == nil {
				f = New()
			}

			defer func() {
				// save flash back to session
				b, err := f.Encode()
				if err == nil {
					sess.Set(flashKey, b)
				}
			}()

			nr := r.WithContext(Set(ctx, f))
			h.ServeHTTP(w, nr)
		})
	}
}

// New creates new flash
func New() Flash {
	return make(Flash)
}

// Decode decodes flash data
func Decode(b []byte) (Flash, error) {
	f := New()
	err := gob.NewDecoder(bytes.NewReader(b)).Decode(&f)
	if err != nil {
		return nil, err
	}
	return f, nil
}

type contextKey int

const flashKey contextKey = iota

// Get gets flash from context
func Get(ctx context.Context) Flash {
	f, ok := ctx.Value(flashKey).(Flash)
	if !ok {
		return New()
	}
	return f
}

// Set sets flash to context's value
func Set(ctx context.Context, f Flash) context.Context {
	return context.WithValue(ctx, flashKey, f)
}

// Encode encodes flash
func (f Flash) Encode() ([]byte, error) {
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
