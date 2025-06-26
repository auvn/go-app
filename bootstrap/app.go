package bootstrap

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/auvn/go-app/graceful"
)

type goFunc = func(ctx context.Context) error
type cleanupFunc = func(err error)

type App struct {
	goFuncs []goFunc
	cleanup []cleanupFunc
}

func (a *App) Go(fns ...func(ctx context.Context) error) {
	a.goFuncs = append(a.goFuncs, fns...)
}

func (a *App) Cleanup(fns ...func(err error)) {
	a.cleanup = append(a.cleanup, fns...)
}

func (a *App) Run(ctx context.Context) {
	eg, ctx := errgroup.WithContext(ctx)

	for _, run := range a.goFuncs {
		run := run
		eg.Go(func() error { return run(ctx) })
	}

	eg.Go(func() error {
		return graceful.ListenSignalContext(ctx)
	})

	err := eg.Wait()

	for _, finalize := range a.cleanup {
		finalize(err)
	}
}
