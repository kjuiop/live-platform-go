package main

import (
	"context"
	"github.com/kjuiop/live-platform-go/cmd/controller/app"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	BUILD_TIME  = "no flag of BUILD_TIME"
	GIT_HASH    = "no flag of GIT_HASH"
	APP_VERSION = "no flag of APP_VERSION"
)

func main() {

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	a := app.NewApplication(ctx, GIT_HASH, APP_VERSION)
	wg.Add(1)
	go a.Start(&wg)

	slog.Info("live chat api app start", "git_hash", GIT_HASH, "build_time", BUILD_TIME, "app_version", APP_VERSION)

	<-exitSignal()
	a.Stop(ctx)
	cancel()
	wg.Wait()
	slog.Info("live chat api app gracefully shutdown")
}

func exitSignal() <-chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	return sig
}
