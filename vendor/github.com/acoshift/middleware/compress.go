package middleware

import (
	"bufio"
	"compress/flate"
	"compress/gzip"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// CompressConfig is the compress middleware config
type CompressConfig struct {
	Skipper   Skipper
	New       func() Compressor
	Encoding  string // http Accept-Encoding, Content-Encoding value
	Vary      bool   // add Vary: Accept-Encoding
	Types     string // only compress for given types, * for all types
	MinLength int    // skip if Content-Length less than given value
}

// default values
const (
	defaultCompressVary      = true
	defaultCompressTypes     = "text/plain text/html text/css text/xml text/javascript application/x-javascript application/xml"
	defaultCompressMinLength = 860
)

// pre-defined compressors
var (
	GzipCompressor = CompressConfig{
		Skipper: DefaultSkipper,
		New: func() Compressor {
			g, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
			if err != nil {
				panic(err)
			}
			return g
		},
		Encoding:  "gzip",
		Vary:      defaultCompressVary,
		Types:     defaultCompressTypes,
		MinLength: defaultCompressMinLength,
	}

	DeflateCompressor = CompressConfig{
		Skipper: DefaultSkipper,
		New: func() Compressor {
			g, err := flate.NewWriter(ioutil.Discard, flate.DefaultCompression)
			if err != nil {
				panic(err)
			}
			return g
		},
		Encoding:  "deflate",
		Vary:      defaultCompressVary,
		Types:     defaultCompressTypes,
		MinLength: defaultCompressMinLength,
	}
)

// Compress creates new compress middleware
func Compress(config CompressConfig) Middleware {
	// fill default config
	if config.Skipper == nil {
		config.Skipper = DefaultSkipper
	}

	mapTypes := make(map[string]struct{})
	for _, t := range strings.Split(config.Types, " ") {
		mapTypes[t] = struct{}{}
	}

	pool := &sync.Pool{
		New: func() interface{} {
			return config.New()
		},
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			// skip if client not support
			if !strings.Contains(r.Header.Get("Accept-Encoding"), config.Encoding) {
				h.ServeHTTP(w, r)
				return
			}

			// skip if web socket
			if r.Header.Get("Sec-WebSocket-Key") != "" {
				h.ServeHTTP(w, r)
				return
			}

			hh := w.Header()

			// skip if already encode
			if hh.Get("Content-Encoding") != "" {
				h.ServeHTTP(w, r)
				return
			}

			if config.Vary {
				addHeaderIfNotExists(hh, "Vary", "Accept-Encoding")
			}

			gw := &compressWriter{
				ResponseWriter: w,
				pool:           pool,
				encoding:       config.Encoding,
				types:          mapTypes,
				minLength:      config.MinLength,
			}
			defer gw.Close()

			h.ServeHTTP(gw, r)
		})
	}
}

// Compressor type
type Compressor interface {
	io.Writer
	io.Closer
	Reset(io.Writer)
	Flush() error
}

type compressWriter struct {
	http.ResponseWriter
	pool        *sync.Pool
	encoder     Compressor
	encoding    string
	types       map[string]struct{}
	wroteHeader bool
	minLength   int
}

func (w *compressWriter) init() {
	h := w.Header()

	// skip if already encode
	if h.Get("Content-Encoding") != "" {
		return
	}

	// skip if length < min length
	if w.minLength > 0 {
		if sl := h.Get("Content-Length"); sl != "" {
			l, _ := strconv.Atoi(sl)
			if l > 0 && l < w.minLength {
				return
			}
		}
	}

	// skip if no match type
	if _, ok := w.types["*"]; !ok {
		ct, _, err := mime.ParseMediaType(h.Get("Content-Type"))
		if err != nil {
			ct = "application/octet-stream"
		}
		if _, ok := w.types[ct]; !ok {
			return
		}
	}

	w.encoder = w.pool.Get().(Compressor)
	w.encoder.Reset(w.ResponseWriter)
	h.Del("Content-Length")
	h.Set("Content-Encoding", w.encoding)
}

func (w *compressWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	if w.encoder != nil {
		return w.encoder.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func (w *compressWriter) Close() {
	if w.encoder == nil {
		return
	}
	w.encoder.Close()
	w.pool.Put(w.encoder)
}

func (w *compressWriter) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.wroteHeader = true
	w.init()
	w.ResponseWriter.WriteHeader(code)
}

// Push implements Pusher interface
func (w *compressWriter) Push(target string, opts *http.PushOptions) error {
	if w, ok := w.ResponseWriter.(http.Pusher); ok {
		return w.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Flush implements Flusher interface
func (w *compressWriter) Flush() {
	if w.encoder != nil {
		w.encoder.Flush()
	}
	if w, ok := w.ResponseWriter.(http.Flusher); ok {
		w.Flush()
	}
}

// CloseNotify implements CloseNotifier interface
func (w *compressWriter) CloseNotify() <-chan bool {
	if w, ok := w.ResponseWriter.(http.CloseNotifier); ok {
		return w.CloseNotify()
	}
	return nil
}

// Hijack implements Hijacker interface
func (w *compressWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w, ok := w.ResponseWriter.(http.Hijacker); ok {
		return w.Hijack()
	}
	return nil, nil, http.ErrNotSupported
}
