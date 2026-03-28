package http

import (
	"net/http"

	"github.com/kjuiop/live-platform-go/internal/shared/models"
	"github.com/kjuiop/live-platform-go/internal/system/adapter/in/http/form"

	"github.com/gin-gonic/gin"

	sysin "github.com/kjuiop/live-platform-go/internal/system/port/in"
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
	group := r.Group("/system")
	group.GET("/health", h.healthCheck)
}

func (s *SystemHandler) successResponse(c *gin.Context, statusCode int, data interface{}) {

	c.JSON(statusCode, models.APIResponse{
		ErrorCode: models.NoError,
		Message:   models.GetCustomMessage(models.NoError),
		Result:    data,
	})
}

func (s *SystemHandler) healthCheck(c *gin.Context) {
	status := s.service.HealthCheck(c.Request.Context())

	systemRes := form.HealthCheckResponse{
		Status:    status.Status,
		Version:   status.Version,
		GitHash:   status.GitHash,
		Uptime:    status.Uptime.String(),
		StartedAt: status.StartedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	s.successResponse(c, http.StatusOK, systemRes)
}
