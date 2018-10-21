package dispatcher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
)

// New creates new dispatcher
func New() *Dispatcher {
	return &Dispatcher{}
}

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
type Dispatcher struct {
	handler map[string]Handler

	Logger *log.Logger
}

func rtName(r reflect.Type) string {
	pkg := r.PkgPath()
	name := r.Name()
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}

func nameFromHandler(h Handler) string {
	return rtName(reflect.TypeOf(h).In(1).Elem())
}

func nameFromMessage(msg Message) string {
	t := reflect.TypeOf(msg)
	if t.Kind() != reflect.Ptr {
		return ""
	}
	return rtName(t.Elem())
}

func isHandler(h Handler) bool {
	t := reflect.TypeOf(h)

	if t.Kind() != reflect.Func {
		return false
	}

	if t.NumIn() != 2 {
		return false
	}
	if t.In(0).Kind() != reflect.Interface && rtName(t.In(0)) != "context.Context" {
		return false
	}
	if t.In(1).Kind() != reflect.Ptr {
		return false
	}

	if t.NumOut() != 1 {
		return false
	}
	if rtName(t.Out(0)) != "error" {
		return false
	}

	return true
}

func (d *Dispatcher) logf(format string, v ...interface{}) {
	if d.Logger != nil {
		d.Logger.Printf(format, v...)
	}
}

// Register registers a handler, and override old handler if exists
func (d *Dispatcher) Register(h Handler) {
	if !isHandler(h) {
		panic("dispatcher: h is not a handler")
	}

	if d.handler == nil {
		d.handler = make(map[string]Handler)
	}

	k := nameFromHandler(h)
	d.handler[k] = h
	d.logf("dispatcher: register %s", k)
}

// Handler returns handler for given message
func (d *Dispatcher) Handler(msg Message) Handler {
	return d.handler[nameFromMessage(msg)]
}

// Dispatch calls handler for given event message
func (d *Dispatcher) Dispatch(ctx context.Context, msg Message) error {
	k := nameFromMessage(msg)
	if k == "" {
		return fmt.Errorf("dispatcher: invalid message type '%s'", reflect.TypeOf(msg))
	}

	d.logf("dispatcher: dispatching %s", k)

	h := d.handler[k]
	if h == nil {
		return ErrNotFound
	}

	err := reflect.ValueOf(h).Call([]reflect.Value{
		reflect.ValueOf(ctx),
		reflect.ValueOf(msg),
	})[0].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}
