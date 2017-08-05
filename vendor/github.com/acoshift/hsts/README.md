# hsts

[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/hsts)](https://goreportcard.com/report/github.com/acoshift/hsts)
[![GoDoc](https://godoc.org/github.com/acoshift/hsts?status.svg)](https://godoc.org/github.com/acoshift/hsts)

HSTS middleware for Golang net/http

### Example

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/acoshift/hsts"
	"github.com/acoshift/middleware"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	h := middleware.Chain(
		hsts.New(hsts.Config{
			MaxAge:            31536000 * time.Second,
			IncludeSubDomains: true,
			Preload:           true,
		}),
	)(http.HandlerFunc(handler))
	http.ListenAndServe(":8080", h)
}
```
