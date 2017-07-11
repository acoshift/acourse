# session

[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/session)](https://goreportcard.com/report/github.com/acoshift/session)
[![GoDoc](https://godoc.org/github.com/acoshift/session?status.svg)](https://godoc.org/github.com/acoshift/session)

Session Middleware for Golang

## Example

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/acoshift/session"
	store "github.com/acoshift/session/store/memory"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		cnt, _ := s.Get("counter").(int)
		cnt++
		s.Set("counter", cnt)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Couter: %d<br><a href=\"/reset\">Reset</a>", cnt)
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context())
		s.Del("counter")
		http.Redirect(w, r, "/", http.StatusFound)
	})

	h := session.Middleware(session.Config{
		Domain:   "localhost",
		HTTPOnly: true,
		Entropy:  32,
		Secret:   []byte("testsecret1234"),
		MaxAge:   time.Minute,
		Name:     "sess",
		Path:     "/",
		Secure:   session.PreferSecure,
		Store:    store.New(store.Config{}),
	})(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}

```
