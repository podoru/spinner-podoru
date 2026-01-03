package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/podoru/spinner-podoru/internal/infrastructure/logger"
)

const RequestIDKey = "request_id"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		requestID, _ := c.Get(RequestIDKey)

		statusCode := c.Writer.Status()

		fields := map[string]interface{}{
			"request_id": requestID,
			"status":     statusCode,
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"latency":    latency.String(),
			"latency_ms": latency.Milliseconds(),
		}

		if userID, exists := c.Get(UserIDKey); exists {
			fields["user_id"] = userID
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		logEntry := log.WithFields(fields)

		switch {
		case statusCode >= 500:
			logEntry.Error("Server error")
		case statusCode >= 400:
			logEntry.Warn("Client error")
		default:
			logEntry.Info("Request completed")
		}
	}
}

func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get(RequestIDKey)

				log.WithFields(map[string]interface{}{
					"request_id": requestID,
					"error":      err,
					"path":       c.Request.URL.Path,
					"method":     c.Request.Method,
				}).Error("Panic recovered")

				c.AbortWithStatusJSON(500, gin.H{
					"success": false,
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "Internal server error",
					},
				})
			}
		}()

		c.Next()
	}
}
