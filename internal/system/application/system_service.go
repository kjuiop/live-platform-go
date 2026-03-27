package application

import (
	"context"
	"time"

	"github.com/kjuiop/live-platform-go/internal/system/domain"
	sysin "github.com/kjuiop/live-platform-go/internal/system/port/in"
)

var _ sysin.SystemService = (*SystemServiceImpl)(nil)

type SystemServiceImpl struct {
	version   string
	gitHash   string
	startedAt time.Time
}

func NewSystemService(version, gitHash string) *SystemServiceImpl {
	return &SystemServiceImpl{
		version:   version,
		gitHash:   gitHash,
		startedAt: time.Now(),
	}
}

func (s SystemServiceImpl) HealthCheck(ctx context.Context) domain.HealthStatus {
	return domain.NewHealthStatus(s.version, s.gitHash, s.startedAt)
}
