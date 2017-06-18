package servertiming

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	t           time.Time
	wroteHeader bool
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true

	diff := float64(time.Now().Sub(w.t)) / float64(time.Millisecond)
	w.Header().Set("Server-Timing", fmt.Sprintf(`total=%.2f; "Total Response Time"`, diff))
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(b)
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
