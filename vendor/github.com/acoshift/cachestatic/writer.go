package cachestatic

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	h           http.Header
	cache       *bytes.Buffer
	code        int
	wroteHeader bool
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.cache.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) Header() http.Header {
	if w.h == nil {
		w.h = cloneHeader(w.ResponseWriter.Header())
	}
	return w.h
}

func contains(arr []string, s string) bool {
	for _, v := range arr {
		if v == s {
			return true
		}
	}
	return false
}

func (w *responseWriter) WriteHeader(code int) {
	w.wroteHeader = true
	w.code = code

	// copy our header to real header
	if w.h != nil {
		h := w.ResponseWriter.Header()
		for k, vv := range w.h {
			for _, v := range vv {
				if !contains(h[k], v) {
					h.Add(k, v)
				}
			}
		}
	}

	w.ResponseWriter.WriteHeader(code)
}

// Push implements Pusher interface
func (w *responseWriter) Push(target string, opts *http.PushOptions) error {
	if w, ok := w.ResponseWriter.(http.Pusher); ok {
		return w.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Flush implements Flusher interface
func (w *responseWriter) Flush() {
	if w, ok := w.ResponseWriter.(http.Flusher); ok {
		w.Flush()
	}
}

// CloseNotify implements CloseNotifier interface
func (w *responseWriter) CloseNotify() <-chan bool {
	if w, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return w.CloseNotify()
	}
	return nil
}

// Hijack implements Hijacker interface
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w, ok := w.ResponseWriter.(http.Hijacker); ok {
		return w.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}

// cloneHeader from net/http/header.go
func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}
