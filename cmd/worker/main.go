package main

import (
	"log/slog"
)

var (
	BUILD_TIME  = "no flag of BUILD_TIME"
	GIT_HASH    = "no flag of GIT_HASH"
	APP_VERSION = "no flag of APP_VERSION"
)

func main() {
	slog.Info("live chat worker app start", "git_hash", GIT_HASH, "build_time", BUILD_TIME, "app_version", APP_VERSION)
}
