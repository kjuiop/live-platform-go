package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Start timer
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		requestId := uuid.New().String()
		c.Set("request_id", requestId)

		// Process request
		c.Next()

		latency := time.Since(start)
		if latency > time.Minute {
			latency = latency.Truncate(time.Second)
		}

		if raw != "" {
			path = path + "?" + raw
		}

		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		logger := slog.With(
			"request_id", requestId,
			"client_ip", c.ClientIP(),
			"method", c.Request.Method,
			"status_code", statusCode,
			"body_size", c.Writer.Size(),
			"path", path,
			"user_agent", c.Request.UserAgent(),
			"latency", latency.String(),
		)

		switch {
		case IsSuccess(statusCode):
			logger.Info("success")
		case IsClientError(statusCode):
			logger.Warn("client error", "error_message", errorMessage)
		case IsServerError(statusCode):
			logger.Error("server error", "error_message", errorMessage)
		default:
			logger.Warn("unexpected status code")
		}
	}
}
