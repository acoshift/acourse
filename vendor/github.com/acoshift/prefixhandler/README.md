# prefixhandler

[![Build Status](https://travis-ci.org/acoshift/prefixhandler.svg?branch=master)](https://travis-ci.org/acoshift/prefixhandler)
[![Coverage Status](https://coveralls.io/repos/github/acoshift/prefixhandler/badge.svg?branch=master)](https://coveralls.io/github/acoshift/prefixhandler?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/prefixhandler)](https://goreportcard.com/report/github.com/acoshift/prefixhandler)
[![GoDoc](https://godoc.org/github.com/acoshift/prefixhandler?status.svg)](https://godoc.org/github.com/acoshift/prefixhandler)

Path Prefix Handler for http.ServeMux

## Example

```go
package main

import (
    "io"
    "net/http"

    "github.com/acoshift/prefixhandler"
)

func main() {
    mux := http.NewServeMux()

    itemMux := http.NewServeMux()
    itemMux.Handle("/", http.HandlerFunc(itemDetail))
    itemMux.Handle("/edit", http.HandlerFunc(itemEdit))

    mux.Handle("/item", http.NotFoundHandler())
    mux.Handle("/item/", prefixhandler.New("/item", "item_id", itemMux))

    http.ListenAndServe(":8080", mux)
}

func itemDetail(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        http.NotFound(w, r)
        return
    }

    itemID := prefixhandler.Get(r.Context(), "item_id")
    // or
    // itemID = r.Context().Value("item_id").(string)

    if itemID == "" {
        http.NotFound(w, r)
        return
    }

    io.WriteString(w, "Item: "+itemID)
}

func itemEdit(w http.ResponseWriter, r *http.Request) {
    itemID := prefixhandler.Get(r.Context(), "item_id")
    io.WriteString(w, "Editing Item: "+itemID)
}
```
