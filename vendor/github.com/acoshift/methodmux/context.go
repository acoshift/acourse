package methodmux

import "context"

type muxKey struct{}

// GetMux gets mux from request's context
// only for fallback handler
func GetMux(ctx context.Context) Mux {
	m, _ := ctx.Value(muxKey{}).(Mux)
	return m
}
