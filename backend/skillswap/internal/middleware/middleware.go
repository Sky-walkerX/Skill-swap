package middleware

import (
	"github.com/gin-gonic/gin"
)

// Backward compatibility functions that use the new organized middleware

// Logger returns the default logging middleware
func Logger() gin.HandlerFunc {
	return RequestLogger()
}

// Recovery returns the default recovery middleware
func Recovery() gin.HandlerFunc {
	return ErrorRecovery()
}

// CORS returns the default CORS middleware
func CORS() gin.HandlerFunc {
	return ConfigurableCORS()
}

// DefaultMiddleware returns a slice of commonly used middleware
func DefaultMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		RequestLogger(),
		ErrorRecovery(),
		ConfigurableCORS(),
		SecurityHeaders(),
	}
}

// ProductionMiddleware returns middleware suitable for production
func ProductionMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		RequestLogger(),
		ErrorRecovery(),
		RateLimit(), // Global rate limiting
		ConfigurableCORS(),
		SecurityHeaders(),
	}
}
