package app

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/kjuiop/live-platform-go/config"
	syscontroller "github.com/kjuiop/live-platform-go/internal/system/adapter/in/http"
	sysapp "github.com/kjuiop/live-platform-go/internal/system/application"
	"github.com/kjuiop/live-platform-go/logger"
	"github.com/kjuiop/live-platform-go/platform/http"
)

type App struct {
	gitHash string
	version string

	cfg *config.EnvConfig
	srv *http.Gin
}

func NewApplication(ctx context.Context, gitHash, version string) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	srv := http.NewGinServer(cfg.Server)

	app := &App{
		gitHash: gitHash,
		version: version,

		cfg: cfg,
		srv: srv,
	}

	app.setupRouter()

	return app
}

func (a *App) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	a.srv.Run()
}

func (a *App) Stop() {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	a.srv.Shutdown(shutdownCtx)
}

func (a *App) setupRouter() {

	// application service register
	sysService := sysapp.NewSystemService(a.version, a.gitHash)

	// controller register
	sysHandler := syscontroller.NewSystemHandler(sysService)

	// router
	router := a.srv.GetEngine()

	// router register
	v1 := router.Group("/api/v1")
	{
		sysHandler.RegisterRoutes(v1.Group("/system"))
	}
}
