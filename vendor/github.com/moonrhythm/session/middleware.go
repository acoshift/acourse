package session

import (
	"bufio"
	"context"
	"errors"
	"net"
	"net/http"
)

// Errors
var (
	ErrNotPassMiddleware = errors.New("session: request not pass middleware")
)

// Middleware is the Manager middleware wrapper
//
// New(config).Middleware()
func Middleware(config Config) func(http.Handler) http.Handler {
	return New(config).Middleware()
}

// Middleware injects session manager into request's context.
//
// All data changed before write response writer's header will be save.
func (m *Manager) Middleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rm := &scopedManager{
				Manager:        m,
				ResponseWriter: w,
				r:              r,
				storage:        make(map[string]*Session),
			}

			ctx := context.WithValue(r.Context(), scopedManagerKey{}, rm)
			h.ServeHTTP(rm, r.WithContext(ctx))

			rm.MustSaveAll()
		})
	}
}

// Get gets session from context
func Get(ctx context.Context, name string) (*Session, error) {
	m, _ := ctx.Value(scopedManagerKey{}).(*scopedManager)
	if m == nil {
		return nil, ErrNotPassMiddleware
	}

	// try get session from storage first
	// to preserve session data from difference handler
	if s, ok := m.storage[name]; ok {
		return s, nil
	}

	// get session from manager
	s, err := m.Get(name)
	if err != nil {
		return nil, err
	}
	s.m = m

	// save session to storage for later get
	m.storage[name] = s
	return s, nil
}

type scopedManagerKey struct{}

type scopedManager struct {
	*Manager
	http.ResponseWriter

	r           *http.Request
	storage     map[string]*Session
	wroteHeader bool
}

func (m *scopedManager) Get(name string) (*Session, error) {
	return m.Manager.Get(m.r, name)
}

func (m *scopedManager) Save(s *Session) error {
	return m.Manager.Save(m.ResponseWriter, s)
}

func (m *scopedManager) MustSaveAll() {
	if m.wroteHeader {
		return
	}

	for _, s := range m.storage {
		err := m.Save(s)
		if err != nil {
			panic("session: " + err.Error())
		}
	}
}

func (m *scopedManager) Regenerate(s *Session) error {
	return m.Manager.Regenerate(m.ResponseWriter, s)
}

func (m *scopedManager) Renew(s *Session) error {
	return m.Manager.Renew(m.ResponseWriter, s)
}

// Write implements http.ResponseWriter
func (m *scopedManager) Write(b []byte) (int, error) {
	if !m.wroteHeader {
		m.WriteHeader(http.StatusOK)
	}
	return m.ResponseWriter.Write(b)
}

// WriteHeader implements http.ResponseWriter
func (m *scopedManager) WriteHeader(code int) {
	if m.wroteHeader {
		return
	}
	m.MustSaveAll()
	m.wroteHeader = true
	m.ResponseWriter.WriteHeader(code)
}

// Push implements Pusher interface
func (m *scopedManager) Push(target string, opts *http.PushOptions) error {
	if w, ok := m.ResponseWriter.(http.Pusher); ok {
		return w.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Flush implements Flusher interface
func (m *scopedManager) Flush() {
	if w, ok := m.ResponseWriter.(http.Flusher); ok {
		w.Flush()
	}
}

// CloseNotify implements CloseNotifier interface
func (m *scopedManager) CloseNotify() <-chan bool {
	if w, ok := m.ResponseWriter.(http.CloseNotifier); ok {
		return w.CloseNotify()
	}
	return nil
}

// Hijack implements Hijacker interface
func (m *scopedManager) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w, ok := m.ResponseWriter.(http.Hijacker); ok {
		return w.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
