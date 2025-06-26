package httpx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type ServeMuxPlugin func(*http.ServeMux)

type Server struct {
	ShutdownContext context.Context
	Addr            string `validate:"nonzero"`
}

func (s *Server) Serve(
	ctx context.Context,
	h http.Handler,
) error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	baseCtx := context.WithoutCancel(ctx)
	baseSrv := http.Server{
		Handler: h,
		BaseContext: func(net.Listener) context.Context {
			return baseCtx
		},
	}

	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	var serveErr error
	go func() {
		defer cancel()
		defer close(done)

		if err := baseSrv.Serve(l); err != nil {
			serveErr = fmt.Errorf("srv.Serve: %w", err)
		}
	}()

	<-ctx.Done()
	shutdownContext := s.ShutdownContext
	if shutdownContext == nil {
		shutdownContext = context.Background()
	}
	shutdownErr := baseSrv.Shutdown(shutdownContext)
	<-done

	return errors.Join(serveErr, shutdownErr)
}

func RunServer(
	ctx context.Context,
	addr string,
	h http.Handler,
) error {
	srv := Server{
		Addr: addr,
	}
	return srv.Serve(ctx, h)
}

func RunInstrumentsServer(
	ctx context.Context,
	addr string,
	plugins ...ServeMuxPlugin,
) error {
	mux := http.NewServeMux()
	for _, p := range plugins {
		p(mux)
	}
	srv := Server{
		Addr: addr,
	}
	return srv.Serve(ctx, mux)
}
