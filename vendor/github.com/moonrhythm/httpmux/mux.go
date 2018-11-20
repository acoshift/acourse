package httpmux

import (
	"net/http"
	"path"
	"strings"
)

// Mux type
type Mux struct {
	m          muxer
	prefix     string
	middleware func(http.Handler) http.Handler
}

type muxer interface {
	http.Handler
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}

// New creates new mux
func New() *Mux {
	return &Mux{
		m: http.NewServeMux(),
	}
}

// Handle registers handler into mux
func (m *Mux) Handle(pattern string, handler http.Handler) {
	if m.middleware != nil {
		handler = m.middleware(handler)
	}

	trailingSlash := len(pattern) > 1 && strings.HasSuffix(pattern, "/")
	pattern = path.Join(m.prefix, pattern)
	if trailingSlash {
		pattern += "/"
	}

	m.m.Handle(pattern, handler)
}

// HandleFunc registers handler into mux
func (m *Mux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.Handle(pattern, http.HandlerFunc(handler))
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.m.ServeHTTP(w, r)
}

// Group creates new group mux
func (m *Mux) Group(prefix string, middleware func(http.Handler) http.Handler) *Mux {
	return &Mux{
		m:          m,
		prefix:     prefix,
		middleware: middleware,
	}
}
