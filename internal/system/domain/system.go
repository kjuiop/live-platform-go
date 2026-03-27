package domain

import "time"

type HealthStatus struct {
	Status    string
	Version   string
	GitHash   string
	Uptime    time.Duration
	StartedAt time.Time
}

func NewHealthStatus(version, gitHash string, startedAt time.Time) HealthStatus {
	return HealthStatus{
		Status:    "ok",
		Version:   version,
		GitHash:   gitHash,
		Uptime:    time.Since(startedAt),
		StartedAt: startedAt,
	}
}
