package application_test

import (
	"context"
	"testing"
	"time"

	"github.com/kjuiop/live-platform-go/internal/system/application"
)

func TestSystemServiceImpl_HealthCheck(t *testing.T) {
	before := time.Now()
	svc := application.NewSystemService("1.0.0", "abc1234")

	status := svc.HealthCheck(context.Background())

	if status.Status != "ok" {
		t.Errorf("Status: got '%s', want 'ok'", status.Status)
	}
	if status.Version != "1.0.0" {
		t.Errorf("Version: got '%s', want '1.0.0'", status.Version)
	}
	if status.GitHash != "abc1234" {
		t.Errorf("GitHash: got '%s', want 'abc1234'", status.GitHash)
	}
	if status.StartedAt.Before(before) {
		t.Error("StartedAt: should be after test start time")
	}
	if status.Uptime < 0 {
		t.Errorf("Uptime: should be non-negative, got %s", status.Uptime)
	}
}
