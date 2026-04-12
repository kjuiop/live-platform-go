PROJECT_PATH=$(shell pwd)

BUILD_NUM_FILE=build_num.txt
BUILD_NUM=$$(cat ./build_num.txt)
APP_VERSION=0.0

TARGET_VERSION=$(APP_VERSION).$(BUILD_NUM)
TARGET_DIR=bin

API_MODULE_NAME=live-chat-api
API_OUTPUT=$(PROJECT_PATH)/$(TARGET_DIR)/$(API_MODULE_NAME)
API_MAIN=cmd/controller/main.go

WORKER_MODULE_NAME=live-chat-worker
WORKER_OUTPUT=$(PROJECT_PATH)/$(TARGET_DIR)/$(WORKER_MODULE_NAME)
WORKER_MAIN=cmd/worker/main.go

LDFLAGS=-X main.BUILD_TIME=`date -u '+%Y-%m-%d_%H:%M:%S'`
LDFLAGS+=-X main.APP_VERSION=$(TARGET_VERSION)
LDFLAGS+=-X main.GIT_HASH=`git rev-parse HEAD`
LDFLAGS+=-s -w

LOCAL_PREFIX=github.com/kjuiop/live-platform-go

.PHONY: build api-build worker-build test fmt lint target-version build_num clean redisinsight

api: config fmt lint api-build

worker: config fmt lint worker-build

build: config fmt lint api-build worker-build

config:
	@if [ ! -d $(TARGET_DIR) ]; then mkdir $(TARGET_DIR); fi

api-build:
	go build -ldflags "$(LDFLAGS)" -o $(API_OUTPUT) $(PROJECT_PATH)/$(API_MAIN)
	cp $(API_OUTPUT) ./live-chat-api

worker-build:
	go build -ldflags "$(LDFLAGS)" -o $(WORKER_OUTPUT) $(PROJECT_PATH)/$(WORKER_MAIN)
	cp $(WORKER_OUTPUT) ./live-chat-worker

test:
	@echo "Running tests..."
	@go clean -testcache
	@go test -race -coverprofile=coverage.out ./...

fmt:
	@echo "Running goimports..."
	@goimports -w -local $(LOCAL_PREFIX) $(shell find . -name "*.go" -not -path "./.git/*")

lint:
	@echo "Running linters..."
	@golangci-lint run --timeout=5m

target-version:
	@echo "========================================"
	@echo "APP_VERSION    : $(APP_VERSION)"
	@echo "BUILD_NUM      : $(BUILD_NUM)"
	@echo "TARGET_VERSION : $(TARGET_VERSION)"
	@echo "========================================"

build_num:
	@echo $$(($$(cat $(BUILD_NUM_FILE)) + 1 )) > $(BUILD_NUM_FILE)
	@echo "BUILD_NUM      : $(BUILD_NUM)"

git-setup: git-template git-hooks
	@echo "✅ Done. (repo-local git template + hooks applied)"

git-template:
	@echo "Setting git commit template..."
	@git config commit.template .gitmessage.txt
	@echo "Done."

git-hooks:
	@echo "Enabling repo hooks (.githooks)..."
	@git config core.hooksPath .githooks
	@chmod +x .githooks/commit-msg
	@chmod +x .githooks/pre-commit
	@echo "Done. (commit-msg & pre-commit hook active)"

redisinsight:
	docker compose --profile dev-tools up redisinsight -d

clean:
	@echo "Cleaning up..."
	@rm -f coverage.out live-chat-api live-chat-worker $(TARGET_DIR)/*
	@echo "Cleanup completed."