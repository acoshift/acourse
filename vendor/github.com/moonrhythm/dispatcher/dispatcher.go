package dispatcher

import (
	"context"
	"errors"
	"reflect"
	"time"
)

// Errors
var (
	ErrNotFound = errors.New("dispatcher: handler not found")
)

// Handler is the event handler
//
// func(context.Context, *Any) error
type Handler interface{}

// Message is the event message
type Message interface{}

// Dispatcher is the event dispatcher
type Dispatcher interface {
	Dispatch(context.Context, Message) error
}

func reflectTypeName(r reflect.Type) string {
	pkg := r.PkgPath()
	name := r.Name()
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

// MessageNameFromHandler gets message name from handler
func MessageNameFromHandler(h Handler) string {
	if !isHandler(h) {
		return ""
	}
	return reflectTypeName(reflect.TypeOf(h).In(1).Elem())
}

// MessageName gets message name
func MessageName(msg Message) string {
	t := reflect.TypeOf(msg)
	if t.Kind() != reflect.Ptr {
		return ""
	}
	return reflectTypeName(t.Elem())
}

func isHandler(h Handler) bool {
	t := reflect.TypeOf(h)

	if t.Kind() != reflect.Func {
		return false
	}

	if t.NumIn() != 2 {
		return false
	}
	if t.In(0).Kind() != reflect.Interface && reflectTypeName(t.In(0)) != "context.Context" {
		return false
	}
	if t.In(1).Kind() != reflect.Ptr {
		return false
	}

	if t.NumOut() != 1 {
		return false
	}
	if reflectTypeName(t.Out(0)) != "error" {
		return false
	}

	return true
}

// Dispatch calls handler for given messages in sequence order,
// when a handler returns error, dispatch will stop and return that error
func Dispatch(ctx context.Context, d Dispatcher, msg ...Message) error {
	for _, m := range msg {
		err := d.Dispatch(ctx, m)
		if err != nil {
			return err
		}
	}
	return nil
}

// DispatchAfter calls dispatch after given duration
// or run immediate if duration is negative,
// then call resultFn with return error
func DispatchAfter(ctx context.Context, d Dispatcher, duration time.Duration, resultFn func(err error), msg ...Message) {
	if resultFn == nil {
		resultFn = func(_ error) {}
	}

	go func() {
		select {
		case <-time.After(duration):
			resultFn(Dispatch(ctx, d, msg...))
		case <-ctx.Done():
			resultFn(ctx.Err())
		}
	}()
}

// DispatchAt calls dispatch at given time,
// and will run immediate if time already passed
func DispatchAt(ctx context.Context, d Dispatcher, t time.Time, resultFn func(err error), msg ...Message) {
	DispatchAfter(ctx, d, time.Until(t), resultFn, msg...)
}
