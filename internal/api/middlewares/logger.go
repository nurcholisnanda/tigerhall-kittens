package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggerMiddleware() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(c *gin.Context) {
		startTime := time.Now()

		// Add request-specific context before handling the request
		requestLogger := logger.WithFields(logrus.Fields{
			"request_id": c.GetString("requestID"), // Assuming you have middleware to add this
			"path":       c.Request.URL.Path,
			"method":     c.Request.Method,
		})

		c.Next()

		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// Log request details and errors
		requestLogger.WithFields(logrus.Fields{
			"status":  c.Writer.Status(),
			"latency": latency,
		}).Info("Request completed")

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				requestLogger.Error(err.Err)
			}
		}
	}
}
