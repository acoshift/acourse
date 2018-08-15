package redisctx

import (
	"context"
	"net/http"

	"github.com/acoshift/middleware"
	"github.com/go-redis/redis"
)

type key int

const (
	keyClient key = iota
	keyPrefix
)

// NewClientContext creates new redis client context
func NewClientContext(ctx context.Context, c *redis.Client) context.Context {
	return context.WithValue(ctx, keyClient, c)
}

// NewPrefixContext creates new redis prefix context
func NewPrefixContext(ctx context.Context, prefix string) context.Context {
	return context.WithValue(ctx, keyPrefix, prefix)
}

// Middleware creates new redis middleware
func Middleware(c *redis.Client, prefix string) middleware.Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = NewClientContext(ctx, c)
			ctx = NewPrefixContext(ctx, prefix)
			r = r.WithContext(ctx)
			h.ServeHTTP(w, r)
		})
	}
}

// GetClient gets client from context
func GetClient(ctx context.Context) *redis.Client {
	return ctx.Value(keyClient).(*redis.Client)
}

// GetPrefix gets prefix from context
func GetPrefix(ctx context.Context) string {
	return ctx.Value(keyPrefix).(string)
}
