# middleware

net/http middleware collection

## Middleware

```go
type Middleware func(h http.Handler) http.Handler
```

### Create new middleware

```go
func say(text string) Middleware {
    return func(h http.Handler) http.Handler {
        return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
            fmt.Println(text)
            h.ServeHTTP(w, r)
        })
    }
}
```

## Chaining

Like normal function middleware can chained.

```go
middleware.HSTSPreload()(middleware.NonWWWRedirect()(say("hello")(handler)))
```

Or using `Chain` to create new middleware

```go
newMiddleware := middleware.Chain(
    middleware.HSTSPreload(),
    middleware.NonWWWRedirect(),
    say("hello"),
)
```

Then

```go
handler = newMiddleware(handler)
```

## HSTS

```go
middleware.HSTS(HSTSConfig{
    MaxAge:            3600 * time.Second,
    IncludeSubDomains: true,
    Preload:           false,
})
```

```go
middleware.HSTS(middleware.DefaultHSTS)
```

```go
middleware.HSTS(middleware.PreloadHSTS)
```

```go
middleware.HSTSPreload()
```

## Compressor

```go
middleware.Compress(middleware.CompressConfig{
    New: func() Compressor {
        g, err := gzip.NewWriterLevel(ioutil.Discard, gzip.DefaultCompression)
        if err != nil {
            panic(err)
        }
        return g
    },
    Encoding:  "gzip",
    Vary:      true,
    Types:     "text/plain text/html",
    MinLength: 1000,
})
```

```go
middleware.Compress(middleware.GzipCompressor)
```

```go
middleware.Compress(middleware.DeflateCompressor)
```

### BrCompressor

```go
middleware.Compress(middleware.BrCompressor)
```

```Dockerfile
FROM alpine

RUN apk add --no-cache ca-certificates tzdata

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk --no-cache add brotli

RUN mkdir -p /app
WORKDIR /app

ADD entrypoint ./
ENTRYPOINT ["/app/entrypoint"]
```

or using `acoshift/go-alpine`

```Dockerfile
FROM acoshift/go-alpine

RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories
RUN apk --no-cache add brotli

RUN mkdir -p /app
WORKDIR /app

ADD entrypoint ./
ENTRYPOINT ["/app/entrypoint"]
```

Builder

use `acoshift/gobuilder` or build your own build image

```Dockerfile
FROM gcr.io/cloud-builders/go

RUN apk --no-cache add cmake

RUN git clone https://github.com/google/brotli && cd brotli && cmake . && make install && cd .. && rm -rf brotli
```

and add `-tags=cbrotli` when using `go build`

### Compress Order

Unlike normal middleware, compressor have to peek on response header.
Order will reverse ex.

```go
middleware.Chain(
    middleware.Compress(middleware.DeflateCompressor),
    middleware.Compress(middleware.GzipCompressor),
    middleware.Compress(middleware.BrCompressor),
)
```

Code above will run `br` first, if client not support `br`, the `gzip` compressor will be run.
Then if client not support both `br` and `gzip`, the `deflate` compressor will be run.

## Redirector

### Redirect from www to non-www

```go
middleware.NonWWWRedirect()
```

### Redirect from non-www to www

```go
middleware.WWWRedirect()
```

### Redirect from http to https

> TODO

## CORS

```go
middleware.CORS(middleware.DefaultCORS) // for public api
```

```go
middleware.CORS(CORSConfig{
    AllowOrigins: []string{"example.com"},
    AllowMethods: []string{
        http.MethodGet,
        http.MethodPost,
    },
    AllowHeaders: []string{
        "Content-Type",
    },
    AllowCredentials: true,
    MaxAge: time.Hour,
})
```

## CSRF

CSRF will reject origin or referal that not in whitelist on `POST`.

```go
middleware.CSRF(middleware.CSRFConfig{
    Origins: []string{
        "http://example.com",
        "https://example.com",
        "http://www.example.com",
        "https://www.example.com",
    },
})
```

or using `IgnoreProto` to ignore protocol

```go
middleware.CSRF(middleware.CSRFConfig{
    Origins: []string{
        "example.com",
        "www.example.com",
    },
    IgnoreProto: true,
})
```

## Ratelimit

> TODO

## Logging

> TODO

## Add Header

`AddHeader` is a helper middleware to add a header to response writer
if not exists

```go
middleware.AddHeader("Vary", "Origin")
```

## License

MIT
