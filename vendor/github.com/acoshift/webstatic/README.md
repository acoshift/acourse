# webstatic

[![Go Report Card](https://goreportcard.com/badge/github.com/acoshift/webstatic)](https://goreportcard.com/report/github.com/acoshift/webstatic)
[![GoDoc](https://godoc.org/github.com/acoshift/webstatic?status.svg)](https://godoc.org/github.com/acoshift/webstatic)

Web Static is the Go handler for handle static files,
returns not found for directory

## Usage

```go
http.Handle("/-/", http.StripPrefix("/-", webstatic.NewDir("assets")))
```

or

```go
http.Handle("/-/", http.StripPrefix("/-", webstatic.New(webstatic.Config{
	Dir: "assets",
	CacheControl: "public, max-age=3600",
})))
```

## License

MIT
