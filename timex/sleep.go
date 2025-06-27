package timex

import (
	"context"
	"time"
)

func Sleep(
	ctx context.Context,
	dur time.Duration,
) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(dur):
		return nil
	}
}
