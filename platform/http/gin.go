package http

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kjuiop/live-platform-go/config"
)

type Gin struct {
	srv    *http.Server
	router *gin.Engine
	cfg    config.Server
}

func NewGinServer(cfg config.Server) (*Gin, error) {

	router := getGinEngine(cfg.Mode)

	proxies := splitAndTrim(cfg.TrustedProxies, ",")
	if err := router.SetTrustedProxies(proxies); err != nil {
		return nil, fmt.Errorf("failed to set trusted proxies: %w", err)
	}

	router.Use(gin.Recovery())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Gin{
		srv:    srv,
		router: router,
		cfg:    cfg,
	}, nil
}

func (g *Gin) Run() {
	err := g.srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		slog.Debug("server close")
	} else {
		log.Fatalf("run server error : %v", err)
	}
}

func (g *Gin) Shutdown(ctx context.Context) {
	if err := g.srv.Shutdown(ctx); err != nil {
		slog.Error("error during server shutdown", "error", err)
	}
}

func (g *Gin) GetEngine() *gin.Engine {
	return g.router
}

func getGinEngine(mode string) *gin.Engine {
	switch mode {
	case "prod":
		return gin.New()
	case "test":
		gin.SetMode(gin.TestMode)
		return gin.Default()
	default:
		return gin.Default()
	}
}

func splitAndTrim(str, sep string) []string {
	parts := strings.Split(str, sep)
	var proxies []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			proxies = append(proxies, trimmed)
		}
	}
	return proxies
}
