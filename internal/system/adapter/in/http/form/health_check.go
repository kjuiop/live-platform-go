package form

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	GitHash   string `json:"git_hash"`
	Uptime    string `json:"uptime"`
	StartedAt string `json:"started_at"`
}
