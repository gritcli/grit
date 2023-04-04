package signalx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
)

// NotifyContextWithCause is a variant of signal.NotifyContextWithCause() that makes the
// signal-related cause of the context's cancellation available via
// SignalCause().
func NotifyContextWithCause(
	ctx context.Context,
	signals ...os.Signal,
) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancelCause(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)

	go func() {
		select {
		case sig := <-c:
			cancel(signalError{sig})
		case <-ctx.Done():
		}
	}()

	return ctx, func() { cancel(nil) }
}

// SignalCause returns the signal that caused the context to be canceled, if
// any.
//
// It returns nil if the context was canceled for any other reason.
func SignalCause(ctx context.Context) os.Signal {
	var err signalError

	if errors.As(context.Cause(ctx), &err) {
		return err.Signal
	}

	return nil
}

type signalError struct {
	Signal os.Signal
}

func (e signalError) Error() string {
	return fmt.Sprintf("received %s signal", e.Signal)
}
