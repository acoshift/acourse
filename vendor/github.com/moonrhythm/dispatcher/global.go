package dispatcher

import "context"

var defaultDispatcher = New()

// Register registers a handler into default dispatcher
func Register(h Handler) {
	defaultDispatcher.Register(h)
}

// Dispatch dispatchs default dispatcher
func Dispatch(ctx context.Context, msg Message) error {
	return defaultDispatcher.Dispatch(ctx, msg)
}
