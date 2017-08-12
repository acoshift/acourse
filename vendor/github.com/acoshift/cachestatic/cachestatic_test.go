package cachestatic

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	mgzip "github.com/acoshift/gzip"
	"github.com/acoshift/header"
	"github.com/acoshift/middleware"
)

func createTestHandler() http.Handler {
	i := 0
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Set(header.ContentType, "text/plain; charset=utf-8")
			w.Header().Set("Custom-Header", "0")
			w.WriteHeader(200)
			return
		}
		if i == 0 {
			i++
			w.Header().Set(header.ContentType, "text/plain; charset=utf-8")
			w.Header().Set("Custom-Header", "0")
			w.WriteHeader(200)
			w.Write([]byte("OK"))
			return
		}
		w.Header().Set(header.ContentType, "text/plain; charset=utf-8")
		w.Header().Set("Custom-Header", "1")
		w.WriteHeader(200)
		w.Write([]byte("Not first response"))
	})
}

func createStaticHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			w.Header().Set(header.ContentType, "text/plain; charset=utf-8")
			w.Header().Set("Custom-Header", "0")
			w.WriteHeader(200)
			return
		}
		w.Header().Set(header.ContentType, "text/plain; charset=utf-8")
		w.Header().Set("Custom-Header", "0")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
}

func TestCachestatic(t *testing.T) {
	ts := httptest.NewServer(New(DefaultConfig)(createTestHandler()))
	defer ts.Close()

	hit := false

	verify := func(resp *http.Response, err error) {
		if err != nil {
			t.Fatalf("expected error to be nil; got %v", err)
		}
		xCache := resp.Header.Get("X-Cache")
		if !hit && xCache != "MISS" {
			t.Fatalf("expected X-Cache to be MISS; got %s", xCache)
		} else if hit && xCache != "HIT" {
			t.Fatalf("expected X-Cache to be HIT; got %s", xCache)
		}
		hit = true
		if resp.Header.Get(header.ContentType) != "text/plain; charset=utf-8" {
			t.Fatalf("invalid Content-Type; got %v", resp.Header.Get(header.ContentType))
		}
		if resp.Header.Get("Custom-Header") != "0" {
			t.Fatalf("invalid Custom-Header; got %v", resp.Header.Get("Custom-Header"))
		}
		defer resp.Body.Close()
		r, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("read response body error; %v", err)
		}
		if bytes.Compare(r, []byte("OK")) != 0 {
			t.Fatalf("invalid response body; got %v", string(r))
		}
	}

	verify(http.Get(ts.URL))
	verify(http.Get(ts.URL))
	verify(http.Get(ts.URL))
	verify(http.Get(ts.URL))
}

func TestWithGzip(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	wg := &sync.WaitGroup{}

	verify := func(resp *http.Response, err error) {
		defer wg.Done()
		if err != nil {
			t.Fatalf("expected error to be nil; got %v", err)
		}
		if resp.Header.Get(header.ContentType) != "text/plain; charset=utf-8" {
			t.Fatalf("invalid Content-Type; got %v", resp.Header.Get(header.ContentType))
		}
		if resp.Header.Get("Custom-Header") != "0" {
			t.Fatalf("invalid Custom-Header; got %v", resp.Header.Get("Custom-Header"))
		}
		defer resp.Body.Close()
		if resp.Request.Method == http.MethodHead {
			return
		}
		if resp.Header.Get(header.ContentEncoding) == header.EncodingGzip && resp.Request.Header.Get(header.AcceptEncoding) != header.EncodingGzip {
			t.Fatalf("request non gzip; got gzip response")
		}
		var body io.Reader
		if resp.Header.Get(header.ContentEncoding) == header.EncodingGzip {
			body, _ = gzip.NewReader(resp.Body)
		} else {
			body = resp.Body
		}
		r, err := ioutil.ReadAll(body)
		if err != nil {
			t.Fatalf("read response body error; %v", err)
		}
		if bytes.Compare(r, []byte("OK")) != 0 {
			t.Fatalf("invalid response body; got %v", string(r))
		}
	}

	var h http.Handler

	l := 100
	run := func() {
		ts := httptest.NewServer(h)
		defer ts.Close()
		wg.Add(l)
		for i := 0; i < l; i++ {
			req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
			if i%2 == 0 {
				req.Header.Set(header.AcceptEncoding, header.EncodingGzip)
			}
			if i%3 == 0 {
				req.Method = http.MethodHead
			}
			go verify(client.Do(req))
		}
		wg.Wait()
	}

	// default config
	h = middleware.Chain(
		mgzip.New(mgzip.Config{Level: mgzip.BestSpeed}),
		New(DefaultConfig),
	)(createTestHandler())
	run()

	// with skip gzip
	h = middleware.Chain(
		New(Config{
			Skipper: func(r *http.Request) bool {
				return !strings.Contains(r.Header.Get(header.AcceptEncoding), header.EncodingGzip)
			},
		}),
		mgzip.New(mgzip.Config{Level: mgzip.BestSpeed}),
	)(createStaticHandler())
	run()

	// with index gzip
	h = middleware.Chain(
		New(Config{
			Indexer: EncodingIndexer(header.EncodingGzip),
		}),
		mgzip.New(mgzip.Config{Level: mgzip.BestSpeed}),
	)(createStaticHandler())
	run()
}

func TestInvalidate(t *testing.T) {
	ch := make(chan interface{})
	ts := httptest.NewServer(New(Config{
		Invalidator: ch,
	})(createTestHandler()))
	defer ts.Close()

	resp, _ := http.Get(ts.URL)
	resp.Body.Close()
	if resp.Header.Get("Custom-Header") != "0" {
		t.Fatalf("Custom-Header must be 0; got %v", resp.Header.Get("Custom-Header"))
	}

	ch <- "GET:/"

	resp, _ = http.Get(ts.URL)
	resp.Body.Close()
	if resp.Header.Get("Custom-Header") != "1" {
		t.Fatalf("Custom-Header must be 1; got %v", resp.Header.Get("Custom-Header"))
	}
}

func TestInvalidateWildcard(t *testing.T) {
	ch := make(chan interface{})
	ts := httptest.NewServer(New(Config{
		Invalidator: ch,
	})(createTestHandler()))
	defer ts.Close()

	resp, _ := http.Get(ts.URL)
	resp.Body.Close()
	if resp.Header.Get("Custom-Header") != "0" {
		t.Fatalf("Custom-Header must be 0; got %v", resp.Header.Get("Custom-Header"))
	}

	ch <- InvalidateAll

	resp, _ = http.Get(ts.URL)
	resp.Body.Close()
	if resp.Header.Get("Custom-Header") != "1" {
		t.Fatalf("Custom-Header must be 1; got %v", resp.Header.Get("Custom-Header"))
	}
}

func TestLastModified(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(header.LastModified, time.Now().UTC().Format(http.TimeFormat))
		fmt.Fprintf(w, "OK")
	})

	ts := httptest.NewServer(New(DefaultConfig)(h))
	defer ts.Close()

	resp, _ := http.Get(ts.URL)
	resp.Body.Close()
	lastModified := resp.Header.Get(header.LastModified)

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set(header.IfModifiedSince, lastModified)

	resp, _ = http.DefaultClient.Do(req)
	resp.Body.Close()
	if resp.StatusCode != http.StatusNotModified {
		t.Fatalf("expected status code to be 304; got %v", resp.StatusCode)
	}
}

func BenchmarkCacheStatic(b *testing.B) {
	ts := httptest.NewServer(New(DefaultConfig)(createTestHandler()))
	defer ts.Close()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL)
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

func BenchmarkNoCacheStatic(b *testing.B) {
	ts := httptest.NewServer(createTestHandler())
	defer ts.Close()
	for i := 0; i < b.N; i++ {
		resp, err := http.Get(ts.URL)
		if err != nil {
			b.Fatal(err)
		}
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}
}

func Example() {
	i := 0
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK %d", i)
		i++
	})

	http.Handle("/", New(DefaultConfig)(h))
	http.ListenAndServe(":8080", nil)
}
