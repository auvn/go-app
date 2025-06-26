package appx

import (
	"context"
	"errors"
	"log/slog"

	"github.com/auvn/go-app/bootstrap"
	"github.com/auvn/go-app/graceful"
)

func NewApp[T any]() *bootstrap.App {
	app := bootstrap.App{}
	app.Cleanup(defaultCleanup())
	return &app
}

func defaultCleanup() func(err error) {
	return func(err error) {
		lvl := slog.LevelError
		if err != nil && errors.Is(err, graceful.ErrShutdown) {
			lvl = slog.LevelWarn
		}

		slog.Log(context.Background(), lvl, "done", "error", err.Error())
	}
}
