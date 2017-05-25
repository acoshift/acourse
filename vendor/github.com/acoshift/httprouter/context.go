package httprouter

import "context"

type paramsKey struct{}

// WithParams adds the params into the context. A modified context is returned.
func WithParams(parent context.Context, ps Params) context.Context {
	return context.WithValue(parent, paramsKey{}, ps)
}

// GetParams gets params from context.
func GetParams(ctx context.Context) Params {
	ps, _ := ctx.Value(paramsKey{}).(Params)
	return ps
}

// GetParam gets a param by name from context
func GetParam(ctx context.Context, name string) string {
	ps := GetParams(ctx)
	if ps == nil {
		return ""
	}
	return ps.ByName(name)
}
