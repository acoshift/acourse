package app

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"cloud.google.com/go/trace"
	"github.com/google/uuid"
)

type responseWriter struct {
	http.ResponseWriter
	header int
}

func (w *responseWriter) WriteHeader(header int) {
	w.header = header
	w.ResponseWriter.WriteHeader(header)
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

// Logger middleware
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		tw := &responseWriter{w, 0}
		ip := r.Header.Get("X-Real-IP")
		if ip == "" {
			ip = r.RemoteAddr
		}
		h.ServeHTTP(tw, r)
		end := time.Now()
		fmt.Printf("%v | %3d | %13v | %s | %s | %s | %s\n",
			end.Format(time.RFC3339),
			tw.header,
			end.Sub(start),
			ip,
			w.Header().Get("X-Request-ID"),
			r.Method,
			path,
		)
	})
}

// Recovery middleware
func Recovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "%v", e)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// RequestID middleware
func RequestID(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-ID", uuid.New().String())
		h.ServeHTTP(w, r)
	})
}

// HSTS middleware
func HSTS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; preload; includeSubDomains")
		h.ServeHTTP(w, r)
	})
}

// Trace middleware
func Trace(client *trace.Client) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			span := client.SpanFromRequest(r)
			defer span.Finish()
			h.ServeHTTP(w, r)
		})
	}
}
