package app

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type App struct {
	srv *http.Server
}

func New(host string, port int, r http.Handler) *App {
	return &App{
		srv: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: r,
		},
	}
}

func (a *App) Start() error {
	if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve http: %w", err)
	}
	return nil
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	return nil
}
