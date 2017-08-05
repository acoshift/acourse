package gzip

import (
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// Copy from compress/gzip
const (
	NoCompression      = gzip.NoCompression
	BestSpeed          = gzip.BestSpeed
	BestCompression    = gzip.BestCompression
	DefaultCompression = gzip.DefaultCompression
	HuffmanOnly        = gzip.HuffmanOnly
)

// Config is the gzip middleware config
type Config struct {
	Skipper middleware.Skipper
	Level   int
}

// DefaultConfig use default compression level
var DefaultConfig = Config{
	Skipper: middleware.DefaultSkipper,
	Level:   DefaultCompression,
}

// New creates new gzip middleware
func New(config Config) middleware.Middleware {
	// fill default config
	if config.Skipper == nil {
		config.Skipper = DefaultConfig.Skipper
	}

	pool := &sync.Pool{
		New: func() interface{} {
			gz, err := gzip.NewWriterLevel(ioutil.Discard, config.Level)
			if err != nil {
				panic(err)
			}
			return gz
		},
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if config.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			if !strings.Contains(r.Header.Get(header.AcceptEncoding), header.EncodingGzip) {
				h.ServeHTTP(w, r)
				return
			}

			if len(r.Header.Get(header.SecWebSocketKey)) > 0 {
				h.ServeHTTP(w, r)
				return
			}

			hh := w.Header()

			if hh.Get(header.ContentEncoding) == header.EncodingGzip {
				h.ServeHTTP(w, r)
				return
			}

			hh.Set(header.Vary, header.AcceptEncoding)

			gw := &responseWriter{
				ResponseWriter: w,
				pool:           pool,
			}
			defer gw.Close()

			h.ServeHTTP(gw, r)
		})
	}
}
