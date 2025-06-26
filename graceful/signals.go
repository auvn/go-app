package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var ErrShutdown = errors.New("graceful: shutdown")

func ListenSignalContext(
	ctx context.Context,
	sig ...os.Signal,
) error {
	sig = append(sig, syscall.SIGTERM, syscall.SIGINT)
	signalCtx, stop := signal.NotifyContext(ctx, sig...)
	defer stop()

	select {
	case <-signalCtx.Done():
		if ctx.Err() != nil {
			return nil
		}
		return errors.Join(signalCtx.Err(), ErrShutdown)
	}
}
