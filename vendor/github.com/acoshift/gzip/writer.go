package gzip

import (
	"bufio"
	"compress/gzip"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/acoshift/header"
)

type responseWriter struct {
	http.ResponseWriter
	pool *sync.Pool
	g    *gzip.Writer
	l    int
}

func (w *responseWriter) init() {
	h := w.Header()
	if len(h.Get(header.ContentEncoding)) > 0 {
		return
	}
	if w.l == 0 {
		if l := h.Get(header.ContentLength); len(l) > 0 {
			w.l, _ = strconv.Atoi(l)
		}
	}
	if w.l > 0 && w.l <= 860 {
		return
	}

	w.g = w.pool.Get().(*gzip.Writer)
	w.g.Reset(w.ResponseWriter)
	h.Del(header.ContentLength)
	h.Set(header.ContentEncoding, header.EncodingGzip)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if w.g == nil {
		w.init()
	}
	if w.g != nil {
		if len(w.Header().Get(header.ContentType)) == 0 {
			w.Header().Set(header.ContentType, http.DetectContentType(b))
		}
		return w.g.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) Close() {
	if w.g == nil {
		return
	}
	w.g.Close()
	w.pool.Put(w.g)
}

func (w *responseWriter) WriteHeader(code int) {
	if w.g == nil {
		w.init()
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
	if w.g != nil {
		w.g.Flush()
	}
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
