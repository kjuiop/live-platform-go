package domain_test

import (
	"testing"
	"time"

	"github.com/kjuiop/live-platform-go/internal/system/domain"
)

func TestNewHealthStatus(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		gitHash    string
		wantStatus string
	}{
		{
			name:       "정상 케이스",
			version:    "1.0.0",
			gitHash:    "abc1234",
			wantStatus: "ok",
		},
		{
			name:       "버전 없음",
			version:    "",
			gitHash:    "abc1234",
			wantStatus: "ok",
		},
		{
			name:       "gitHash 없음",
			version:    "1.0.0",
			gitHash:    "",
			wantStatus: "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startedAt := time.Now().Add(-5 * time.Second)

			status := domain.NewHealthStatus(tt.version, tt.gitHash, startedAt)

			if status.Status != tt.wantStatus {
				t.Errorf("Status: got '%s', want '%s'", status.Status, tt.wantStatus)
			}
			if status.Version != tt.version {
				t.Errorf("Version: got '%s', want '%s'", status.Version, tt.version)
			}
			if status.GitHash != tt.gitHash {
				t.Errorf("GitHash: got '%s', want '%s'", status.GitHash, tt.gitHash)
			}
			if status.Uptime < 5*time.Second {
				t.Errorf("Uptime: got %s, want >= 5s", status.Uptime)
			}
			if status.StartedAt != startedAt {
				t.Error("StartedAt: does not match input")
			}
		})
	}
}
