package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/availability"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupAvailabilityRoutes sets up all availability-related routes
func SetupAvailabilityRoutes(api *gin.RouterGroup, cfg *config.Config, availabilityHandler *availability.Handler) {
	// All availability routes require authentication
	availabilityGroup := api.Group("/availability")
	availabilityGroup.Use(middleware.JWTAuth(*cfg))
	{
		availabilityGroup.POST("", availabilityHandler.CreateAvailabilitySlot)       // POST /api/v1/availability
		availabilityGroup.GET("", availabilityHandler.GetUserAvailabilitySlots)      // GET /api/v1/availability
		availabilityGroup.GET("/:id", availabilityHandler.GetAvailabilitySlot)       // GET /api/v1/availability/:id
		availabilityGroup.PUT("/:id", availabilityHandler.UpdateAvailabilitySlot)    // PUT /api/v1/availability/:id
		availabilityGroup.DELETE("/:id", availabilityHandler.DeleteAvailabilitySlot) // DELETE /api/v1/availability/:id

		// Advanced availability queries
		availabilityGroup.GET("/search", availabilityHandler.GetAvailabilityByDayAndTime)     // GET /api/v1/availability/search
		availabilityGroup.GET("/common/:user_id", availabilityHandler.FindCommonAvailability) // GET /api/v1/availability/common/:user_id
	}
}
