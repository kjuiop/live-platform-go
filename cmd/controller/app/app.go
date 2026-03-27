package app

import (
	"context"
	"github.com/kjuiop/live-platform-go/config"
	"github.com/kjuiop/live-platform-go/logger"
	"log"
	"sync"
)

type App struct {
	cfg *config.EnvConfig
}

func NewApplication(ctx context.Context) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	return &App{
		cfg: cfg,
	}
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()
}

func (a *App) Stop(ctx context.Context) {
}
