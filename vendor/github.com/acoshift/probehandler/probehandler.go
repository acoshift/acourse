package probehandler

import (
	"net/http"
	"sync"
)

// New creates new handler
func New() *Handler {
	return &Handler{}
}

// Handler is the probe handler
// default handler is success until call Fail
type Handler struct {
	http.Handler

	m sync.RWMutex
	f bool
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.m.RLock()
	f := h.f
	h.m.RUnlock()
	if f {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Success marks probe as success
func (h *Handler) Success() {
	h.m.Lock()
	h.f = false
	h.m.Unlock()
}

// Fail marks probe as fail
func (h *Handler) Fail() {
	h.m.Lock()
	h.f = true
	h.m.Unlock()
}

// Success always send http.StatusOK
func Success() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

// Fail always send http.StatusServiceUnavailable
func Fail() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
}
