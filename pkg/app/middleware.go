package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type loggerWriter struct {
	http.ResponseWriter
	header int
}

func (w *loggerWriter) WriteHeader(header int) {
	w.header = header
	w.ResponseWriter.WriteHeader(header)
}

// Logger middleware
func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.URL.Path
		tw := &loggerWriter{w, 0}
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
		w.Header().Add("Strict-Transport-Security", "max-age=63072000") // ; includeSubDomains
		h.ServeHTTP(w, r)
	})
}
