package app

import (
	"context"
	"github.com/kjuiop/live-platform-go/config"
	"github.com/kjuiop/live-platform-go/logger"
	"github.com/kjuiop/live-platform-go/platform/http"
	"log"
	"sync"
)

type App struct {
	cfg *config.EnvConfig
	srv *http.Gin
}

func NewApplication(ctx context.Context) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	srv := http.NewGinServer(cfg.Server)

	return &App{
		cfg: cfg,
		srv: srv,
	}
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	a.srv.Run()
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Shutdown(ctx)
}
