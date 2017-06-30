package cachestatic

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

// Config type
type Config struct {
	Skipper     middleware.Skipper
	Indexer     Indexer
	Invalidator chan interface{}
	SkipHeaders SkipHeaderFunc
}

type invalidateAll struct{}

// InvalidateAll invalidates all cache items
var InvalidateAll interface{} = invalidateAll{}

// SkipHeaderFunc is the function to skip header,
// return true to skip
type SkipHeaderFunc func(string) bool

// SkipHeaders skips all given headers
func SkipHeaders(headers ...string) SkipHeaderFunc {
	skipHeader := make(map[string]struct{})
	for _, h := range headers {
		skipHeader[h] = struct{}{}
	}
	return func(h string) bool {
		_, ok := skipHeader[h]
		return ok
	}
}

// DefaultConfig is the default config
var DefaultConfig = Config{
	Skipper:     middleware.DefaultSkipper,
	Indexer:     DefaultIndexer,
	Invalidator: nil,
	SkipHeaders: SkipHeaders(header.SetCookie),
}

// New creates new cachestatic middleware
func New(c Config) func(http.Handler) http.Handler {
	if c.Skipper == nil {
		c.Skipper = DefaultConfig.Skipper
	}
	if c.Indexer == nil {
		c.Indexer = DefaultConfig.Indexer
	}
	if c.Invalidator == nil {
		c.Invalidator = DefaultConfig.Invalidator
	}
	if c.SkipHeaders == nil {
		c.SkipHeaders = DefaultConfig.SkipHeaders
	}

	var (
		l     = &sync.RWMutex{}
		cache = make(map[interface{}]*item)
	)

	if c.Invalidator != nil {
		go func() {
			for {
				p := <-c.Invalidator
				l.Lock()
				if _, ok := p.(invalidateAll); ok {
					cache = make(map[interface{}]*item)
				} else {
					delete(cache, p)
				}
				l.Unlock()
			}
		}()
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if c.Skipper(r) {
				h.ServeHTTP(w, r)
				return
			}

			p := c.Indexer(r)
			l.RLock()
			if ci := cache[p]; ci != nil {
				l.RUnlock()
				wh := w.Header()
				for k, vs := range ci.header {
					if !c.SkipHeaders(k) {
						wh[k] = vs
					}
				}

				// check Last-Modified
				if !ci.modTime.IsZero() {
					if ts := r.Header.Get(header.IfModifiedSince); len(ts) > 0 {
						t, _ := http.ParseTime(ts)
						if ci.modTime.Equal(t) {
							wh.Del(header.ContentType)
							wh.Del(header.ContentLength)
							wh.Del(header.AcceptRanges)
							w.WriteHeader(http.StatusNotModified)
							return
						}
					}
				}

				io.Copy(w, bytes.NewReader(ci.data))
				return
			}
			l.RUnlock()
			cw := &responseWriter{
				ResponseWriter: w,
				cache:          &bytes.Buffer{},
			}
			h.ServeHTTP(cw, r)

			// cache only status ok
			if cw.code == http.StatusOK {
				l.Lock()
				cache[p] = createItem(cw)
				l.Unlock()
			}
		})
	}
}
