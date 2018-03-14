package methodmux

import (
	"context"
	"net/http"
)

// Mux is the method mux
type Mux map[string]http.Handler

var _ http.Handler = Mux{}

// FallbackHandler is the default fallback handler if no method matched
var FallbackHandler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	m := GetMux(r.Context())
	for method := range m {
		if method != "" {
			w.Header().Add("Allow", method)
		}
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
})

func (m Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := m[r.Method]; h != nil {
		h.ServeHTTP(w, r)
		return
	}

	// Head fallback to get
	if r.Method == http.MethodHead {
		if h := m[http.MethodGet]; h != nil {
			h.ServeHTTP(w, r)
			return
		}
	}

	r = r.WithContext(context.WithValue(r.Context(), muxKey{}, m))

	// fallback
	if h := m[""]; h != nil {
		h.ServeHTTP(w, r)
		return
	}

	// handler not found, fallback to package's not found handler
	FallbackHandler.ServeHTTP(w, r)
}

// Get is a short-hand for Mux{http.MethodGet: h}
func Get(h http.Handler) Mux {
	return Mux{http.MethodGet: h}
}

// Post is a short-hand for Mux{http.MethodPost: h}
func Post(h http.Handler) Mux {
	return Mux{http.MethodPost: h}
}

// Put is a short-hand for Mux{http.MethodPut: h}
func Put(h http.Handler) Mux {
	return Mux{http.MethodPut: h}
}

// Patch is a short-hand for Mux{http.MethodPatch: h}
func Patch(h http.Handler) Mux {
	return Mux{http.MethodPatch: h}
}

// Delete is a short-hand for Mux{http.MethodDelete: h}
func Delete(h http.Handler) Mux {
	return Mux{http.MethodDelete: h}
}

// Head is a short-hand for Mux{http.MethodHead: h}
func Head(h http.Handler) Mux {
	return Mux{http.MethodHead: h}
}

// Options is a short-hand for Mux{http.MethodOptions: h}
func Options(h http.Handler) Mux {
	return Mux{http.MethodOptions: h}
}

// GetPost is a short-hand for Mux{http.MethodGet: get, http.MethodPost: post}
func GetPost(get, post http.Handler) Mux {
	return Mux{http.MethodGet: get, http.MethodPost: post}
}
