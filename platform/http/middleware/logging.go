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

		// Fill the params
		param := gin.LogFormatterParams{}

		param.TimeStamp = time.Now() // Stop timer
		param.Latency = param.TimeStamp.Sub(start)
		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = c.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path

		logger := slog.With(
			"request_id", requestId,
			"client_ip", param.ClientIP,
			"method", param.Method,
			"status_code", param.StatusCode,
			"body_size", param.BodySize,
			"path", param.Path,
			"user_agent", c.Request.UserAgent(),
			"latency", param.Latency.String(),
		)

		switch {
		case IsSuccess(param.StatusCode):
			logger.Info("success")
		case IsClientError(param.StatusCode):
			logger.Warn("client error", "error_message", param.ErrorMessage)
		case IsServerError(param.StatusCode):
			logger.Error("server error", "error_message", param.ErrorMessage)
		default:
			logger.Warn("unexpected status code")
		}
	}
}
