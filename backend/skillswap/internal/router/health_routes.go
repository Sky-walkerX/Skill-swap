package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupHealthRoutes configures health check and system status routes
func SetupHealthRoutes(router *gin.Engine) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "skillswap-api",
			"version": "1.0.0",
		})
	})

	// Readiness probe (can be extended to check database connectivity, etc.)
	router.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	// Liveness probe
	router.GET("/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})
}
