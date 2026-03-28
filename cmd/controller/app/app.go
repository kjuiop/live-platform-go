package app

import (
	"context"
	"log"
	"sync"
	"time"

	roomcontroller "github.com/kjuiop/live-platform-go/internal/room/adapter/in/http"
	roomRepo "github.com/kjuiop/live-platform-go/internal/room/adapter/out/redis"
	roomapp "github.com/kjuiop/live-platform-go/internal/room/application"

	"github.com/kjuiop/live-platform-go/platform/redis"

	"github.com/kjuiop/live-platform-go/config"
	syscontroller "github.com/kjuiop/live-platform-go/internal/system/adapter/in/http"
	sysapp "github.com/kjuiop/live-platform-go/internal/system/application"
	"github.com/kjuiop/live-platform-go/logger"
	"github.com/kjuiop/live-platform-go/platform/http"
)

type App struct {
	gitHash string
	version string

	cfg  *config.EnvConfig
	srv  *http.Gin
	rcli *redis.Client
}

func NewApplication(ctx context.Context, gitHash, version string) *App {

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("failed to load env config: %v", err)
	}

	if err := logger.SlogInit(cfg.Logger); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	rcli, err := redis.NewRedisSingleClient(ctx, cfg.Redis)
	if err != nil {
		log.Fatalf("failed to initialize redis client: %v", err)
	}

	srv, err := http.NewGinServer(cfg.Server)
	if err != nil {
		log.Fatalf("failed to initialize http server: %v", err)
	}

	app := &App{
		gitHash: gitHash,
		version: version,

		cfg:  cfg,
		srv:  srv,
		rcli: rcli,
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
	a.rcli.Close()
}

func (a *App) setupRouter() {

	timeout := time.Duration(a.cfg.Policy.ContextTimeout) * time.Second

	// repository
	roomRepository := roomRepo.NewRoomRedisRepository(a.rcli)

	// application service register
	sysService := sysapp.NewSystemService(a.version, a.gitHash)
	roomService := roomapp.NewRoomService(timeout, roomRepository)

	// controller register
	sysHandler := syscontroller.NewSystemHandler(sysService)
	roomHandler := roomcontroller.NewRoomHandler(a.cfg.Policy, roomService)

	// router
	router := a.srv.GetEngine()

	// router register
	v1 := router.Group("/api/v1")
	{
		sysHandler.RegisterRoutes(v1)
		roomHandler.RegisterRoutes(v1)
	}
}
