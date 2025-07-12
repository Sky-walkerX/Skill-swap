package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/middleware"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/search"
	"github.com/gin-gonic/gin"
)

func SetupSearchRoutes(api *gin.RouterGroup, searchService service.SearchService, cfg *config.Config) {
	searchHandler := search.NewHandler(searchService)
	
	// Public search routes
	public := api.Group("/search")
	{
		public.GET("/global", searchHandler.GlobalSearch)                    // GET /api/search/global
		public.GET("/suggestions", searchHandler.SearchSuggestions)         // GET /api/search/suggestions
	}
	
	// Protected search routes (for more detailed searches)
	protected := api.Group("/search")
	protected.Use(middleware.JWTAuth(*cfg))
	{
		protected.GET("/users", searchHandler.SearchUsers)                  // GET /api/search/users
		protected.GET("/swaps", searchHandler.SearchSwaps)                  // GET /api/search/swaps
		protected.GET("/skills", searchHandler.SearchSkills)               // GET /api/search/skills
	}
}
