package in

import (
	"context"
	"github.com/kjuiop/live-platform-go/internal/system/domain"
)

type SystemService interface {
	HealthCheck(ctx context.Context) domain.HealthStatus
}
