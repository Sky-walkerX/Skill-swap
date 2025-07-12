package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/admin"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/middleware"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/skill"
	"github.com/gin-gonic/gin"
)

// SetupAdminRoutes configures all admin-related routes
func SetupAdminRoutes(api *gin.RouterGroup, cfg *config.Config, skillHandler *skill.Handler, adminHandler *admin.Handler) {
	// Admin routes group with authentication and admin role check
	adminGroup := api.Group("/admin")
	adminGroup.Use(middleware.JWTAuth(*cfg))
	adminGroup.Use(middleware.AdminAuth())
	{
		// Admin skill management
		skills := adminGroup.Group("/skills")
		{
			skills.POST("", skillHandler.CreateSkill)       // POST /api/v1/admin/skills
			skills.PUT("/:id", skillHandler.UpdateSkill)    // PUT /api/v1/admin/skills/:id
			skills.DELETE("/:id", skillHandler.DeleteSkill) // DELETE /api/v1/admin/skills/:id
		}

		// Admin user management
		users := adminGroup.Group("/users")
		{
			users.GET("", adminHandler.GetAllUsers)                      // GET /api/v1/admin/users
			users.PUT("/:id/ban", adminHandler.BanUser)                  // PUT /api/v1/admin/users/:id/ban
			users.PUT("/:id/unban", adminHandler.UnbanUser)              // PUT /api/v1/admin/users/:id/unban
			users.DELETE("/:id", adminHandler.DeleteUser)                // DELETE /api/v1/admin/users/:id
			users.PUT("/:id/make-admin", adminHandler.MakeUserAdmin)     // PUT /api/v1/admin/users/:id/make-admin
			users.PUT("/:id/remove-admin", adminHandler.RemoveUserAdmin) // PUT /api/v1/admin/users/:id/remove-admin
		}

		// Admin swap management
		swaps := adminGroup.Group("/swaps")
		{
			swaps.GET("", adminHandler.GetAllSwaps)           // GET /api/v1/admin/swaps
			swaps.PUT("/:id/cancel", adminHandler.CancelSwap) // PUT /api/v1/admin/swaps/:id/cancel
		}

		// Platform statistics and monitoring
		adminGroup.GET("/stats", adminHandler.GetPlatformStats)     // GET /api/v1/admin/stats
		adminGroup.GET("/reports", adminHandler.GetReportedContent) // GET /api/v1/admin/reports

		// TODO: Additional admin features
		// - GET /admin/audit-logs - View admin action logs
		// - POST /admin/messages/broadcast - Send platform-wide message
		// - GET /admin/skills/pending - Get skills pending approval (if approval system is implemented)
		// - PUT /admin/skills/:id/approve - Approve skill
		// - PUT /admin/skills/:id/reject - Reject skill
	}
}
