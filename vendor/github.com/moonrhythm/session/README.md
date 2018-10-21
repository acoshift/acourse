# session

[![Build Status](https://travis-ci.org/moonrhythm/session.svg?branch=master)](https://travis-ci.org/moonrhythm/session)
[![codecov](https://codecov.io/gh/moonrhythm/session/branch/master/graph/badge.svg)](https://codecov.io/gh/moonrhythm/session)
[![Go Report Card](https://goreportcard.com/badge/github.com/moonrhythm/session)](https://goreportcard.com/report/github.com/moonrhythm/session)
[![GoDoc](https://godoc.org/github.com/moonrhythm/session?status.svg)](https://godoc.org/github.com/moonrhythm/session)

Session Middleware for Golang

## Example with Middleware

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/moonrhythm/session"
    store "github.com/moonrhythm/session/store/memory"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/" {
            http.NotFound(w, r)
            return
        }

        s, _ := session.Get(r.Context(), "sess")
        cnt := s.GetInt("counter")
        cnt++
        s.Set("counter", cnt)
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, "Couter: %d<br><a href=\"/reset\">Reset</a>", cnt)
    })
    mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
        s, _ := session.Get(r.Context(), "sess")
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

    "github.com/moonrhythm/session"
    store "github.com/moonrhythm/session/store/memory"
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

        s, _ := m.Get(r, "sess")
        cnt := s.GetInt("counter")
        cnt++
        s.Set("counter", cnt)
        m.Save(w, s)
        w.Header().Set("Content-Type", "text/html")
        fmt.Fprintf(w, "Couter: %d<br><a href=\"/reset\">Reset</a>", cnt)
    })
    mux.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
        s, _ := m.Get(r, "sess")
        s.Del("counter")
        m.Save(w, s)
        http.Redirect(w, r, "/", http.StatusFound)
    })

    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## License

MIT
