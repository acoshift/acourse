package dispatcher

import (
	"context"
	"fmt"
	"reflect"
)

// Mux is the dispatch multiplexer
type Mux struct {
	handler map[string]Handler
}

// NewMux creates new mux
func NewMux() *Mux {
	return &Mux{}
}

// Register registers handlers, and override old handler if exists
func (d *Mux) Register(hs ...Handler) {
	if d.handler == nil {
		d.handler = make(map[string]Handler)
	}

	for _, h := range hs {
		k := MessageNameFromHandler(h)
		if k == "" {
			panic("dispatcher: h is not a handler")
		}
		d.handler[k] = h
	}
}

// Handler returns handler for given message
func (d *Mux) Handler(msg Message) Handler {
	return d.handler[MessageName(msg)]
}

// Dispatch calls handler for given messages
func (d *Mux) Dispatch(ctx context.Context, msg Message) error {
	k := MessageName(msg)
	if k == "" {
		return fmt.Errorf("dispatcher: invalid message type '%s'", reflect.TypeOf(msg))
	}

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
