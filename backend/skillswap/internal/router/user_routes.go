package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/middleware"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/user"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(api *gin.RouterGroup, userService service.UserService, cfg *config.Config) {
	userHandler := user.NewHandler(userService)

	// Public user routes (no authentication required)
	public := api.Group("/public")
	{
		public.GET("/users/search", userHandler.SearchUsers)
	}

	// Protected user routes (authentication required)
	protected := api.Group("/users")
	protected.Use(middleware.JWTAuth(*cfg))
	{
		protected.GET("/profile", userHandler.GetProfile)
		protected.PUT("/profile", userHandler.UpdateProfile)
	}
}
