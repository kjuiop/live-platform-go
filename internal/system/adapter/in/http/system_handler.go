package http

import (
	"github.com/gin-gonic/gin"
	sysin "github.com/kjuiop/live-platform-go/internal/system/port/in"
	"net/http"
)

type SystemHandler struct {
	service sysin.SystemService
}

func NewSystemHandler(service sysin.SystemService) *SystemHandler {
	return &SystemHandler{
		service: service,
	}
}

func (h *SystemHandler) RegisterRoutes(r gin.IRouter) {
	r.GET("/health", h.healthCheck)
}

type HealthCheckResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	GitHash   string `json:"git_hash"`
	Uptime    string `json:"uptime"`
	StartedAt string `json:"started_at"`
}

func (h *SystemHandler) healthCheck(c *gin.Context) {
	status := h.service.HealthCheck(c.Request.Context())
	c.JSON(http.StatusOK, HealthCheckResponse{
		Status:    status.Status,
		Version:   status.Version,
		GitHash:   status.GitHash,
		Uptime:    status.Uptime.String(),
		StartedAt: status.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}
