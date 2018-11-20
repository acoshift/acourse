package bus

import (
	"context"
	"log"
	"os"

	"github.com/moonrhythm/dispatcher"
)

// Message is the bus event message
type Message dispatcher.Message

var (
	d      = dispatcher.NewMux()
	logger = log.New(os.Stdout, "", log.LstdFlags)
)

// Dispatch dispatchs event to event bus
func Dispatch(ctx context.Context, msg ...Message) error {
	for _, m := range msg {
		logger.Printf(dispatcher.MessageName(m))
		err := d.Dispatch(ctx, m)
		if err != nil {
			return err
		}
	}
	return nil
}

// Register registers event handler
func Register(h ...dispatcher.Handler) {
	d.Register(h...)
}
