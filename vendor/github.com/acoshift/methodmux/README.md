# methodmux

[![Build Status](https://travis-ci.org/acoshift/methodmux.svg?branch=master)](https://travis-ci.org/acoshift/methodmux)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/methodmux/badge.svg?branch=master)](https://coveralls.io/github/acoshift/methodmux?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/methodmux)](https://goreportcard.com/report/github.com/acoshift/methodmux)
[![GoDoc](https://godoc.org/github.com/acoshift/methodmux?status.svg)](https://godoc.org/github.com/acoshift/methodmux)

Method Multiplexer for http.ServeMux

## Example

```go
package main

import (
    "io"
    "net/http"

    "github.com/acoshift/methodmux"
)

func main() {
    mux := http.NewServeMux()
    mux.Handle("/", methodmux.Get(http.HandlerFunc(index)))
    mux.Handle("/about", methodmux.Mux{
        http.MethodGet:  http.HandlerFunc(aboutGet),
        http.MethodPost: http.HandlerFunc(aboutPost),
        "":              http.HandlerFunc(aboutOther),
    })
    http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "Hello, Method Mux!")
}

func aboutGet(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "About Get")
}

func aboutPost(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "About Post")
}

func aboutOther(w http.ResponseWriter, r *http.Request) {
    io.WriteString(w, "About does not support method "+r.Method)
}
```
