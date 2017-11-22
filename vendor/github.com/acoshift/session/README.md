# session

[![Build Status](https://travis-ci.org/acoshift/session.svg?branch=master)](https://travis-ci.org/acoshift/session)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/session/badge.svg?branch=master)](https://coveralls.io/github/acoshift/session?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/session)](https://goreportcard.com/report/github.com/acoshift/session)
[![GoDoc](https://godoc.org/github.com/acoshift/session?status.svg)](https://godoc.org/github.com/acoshift/session)

Session Middleware for Golang

## Example with Middleware

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
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		s := session.Get(r.Context(), "sess")
		cnt, _ := s.Get("counter").(int)
		cnt++
		s.Set("counter", cnt)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Couter: %d<br><a href=\"/reset\">Reset</a>", cnt)
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		s := session.Get(r.Context(), "sess")
		s.Del("counter")
		http.Redirect(w, r, "/", http.StatusFound)
	})

	h := session.Middleware(session.Config{
		Domain:   "localhost",
		HTTPOnly: true,
		Secret:   []byte("testsecret1234"),
		MaxAge:   time.Minute,
		Path:     "/",
		Secure:   session.PreferSecure,
		Store:    store.New(store.Config{}),
	})(mux)
	// equals to
	// h := session.New(session.Config{...}).Middleware()(mux)

	log.Fatal(http.ListenAndServe(":8080", h))
}

```

## Example with Manager

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

	m := session.New(session.Config{
		Domain:   "localhost",
		HTTPOnly: true,
		Secret:   []byte("testsecret1234"),
		MaxAge:   time.Minute,
		Path:     "/",
		Secure:   session.PreferSecure,
		Store:    store.New(store.Config{}),
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		s := m.Get(r, "sess")
		cnt, _ := s.Get("counter").(int)
		cnt++
		s.Set("counter", cnt)
		m.Save(w, s)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "Couter: %d<br><a href=\"/reset\">Reset</a>", cnt)
	})
	mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		s := m.Get(r, "sess")
		s.Del("counter")
		m.Save(w, s)
		http.Redirect(w, r, "/", http.StatusFound)
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}

```
