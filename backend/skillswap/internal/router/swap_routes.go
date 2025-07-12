package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/gin-gonic/gin"
)

// SetupSwapRoutes configures all swap request related routes
func SetupSwapRoutes(api *gin.RouterGroup, cfg *config.Config) {
	// TODO: Implement swap routes
	// - POST /swaps - Create swap request
	// - GET /swaps - Get user's swap requests (both sent and received)
	// - GET /swaps/:id - Get specific swap request
	// - PUT /swaps/:id/accept - Accept swap request
	// - PUT /swaps/:id/reject - Reject swap request
	// - DELETE /swaps/:id - Cancel/Delete swap request
	// - POST /swaps/:id/rating - Rate completed swap

	_ = api // Unused for now
	_ = cfg // Unused for now
}
