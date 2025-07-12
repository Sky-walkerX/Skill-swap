package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig holds configuration for request logging
type LoggerConfig struct {
	EnableColors    bool
	EnableDetails   bool
	SkipPaths       []string
	CustomFormatter func(gin.LogFormatterParams) string
}

// DefaultLoggerConfig returns default logger configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		EnableColors:  false,
		EnableDetails: true,
		SkipPaths:     []string{"/health", "/ready", "/live"},
	}
}

// RequestLogger returns a configurable logging middleware
func RequestLogger(config ...LoggerConfig) gin.HandlerFunc {
	cfg := DefaultLoggerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Use custom formatter if provided, otherwise use default
	if cfg.CustomFormatter != nil {
		return gin.LoggerWithConfig(gin.LoggerConfig{
			Formatter: cfg.CustomFormatter,
			SkipPaths: cfg.SkipPaths,
		})
	}

	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: func(param gin.LogFormatterParams) string {
			if !cfg.EnableDetails {
				return fmt.Sprintf("%s %s %d\n", param.Method, param.Path, param.StatusCode)
			}

			return fmt.Sprintf("[%s] %s - \"%s %s %s\" %d %s \"%s\" %s\n",
				param.TimeStamp.Format(time.RFC3339),
				param.ClientIP,
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
			)
		},
		SkipPaths: cfg.SkipPaths,
	})
}

// ErrorRecovery returns a recovery middleware with custom error handling
func ErrorRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic for debugging
		fmt.Printf("Panic recovered: %v\n", recovered)

		c.JSON(500, gin.H{
			"error":   "Internal server error",
			"message": "Something went wrong. Please try again later.",
		})
	})
}
