package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/middleware"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/skill"
	"github.com/gin-gonic/gin"
)

// SetupSkillRoutes configures all skill-related routes
func SetupSkillRoutes(api *gin.RouterGroup, cfg *config.Config, skillHandler *skill.Handler) {
	// Public skill routes (no authentication required)
	skills := api.Group("/skills")
	{
		skills.GET("", skillHandler.GetAllSkills) // GET /api/v1/skills
		skills.GET("/:id", skillHandler.GetSkill) // GET /api/v1/skills/:id
	}

	// Protected user skill routes (authentication required)
	userSkills := api.Group("/users/skills")
	userSkills.Use(middleware.JWTAuth(*cfg))
	{
		// Offered skills
		userSkills.GET("/offered", skillHandler.GetUserOfferedSkills)      // GET /api/v1/users/skills/offered
		userSkills.POST("/offered", skillHandler.AddOfferedSkill)          // POST /api/v1/users/skills/offered
		userSkills.DELETE("/offered/:id", skillHandler.RemoveOfferedSkill) // DELETE /api/v1/users/skills/offered/:id

		// Wanted skills
		userSkills.GET("/wanted", skillHandler.GetUserWantedSkills)      // GET /api/v1/users/skills/wanted
		userSkills.POST("/wanted", skillHandler.AddWantedSkill)          // POST /api/v1/users/skills/wanted
		userSkills.DELETE("/wanted/:id", skillHandler.RemoveWantedSkill) // DELETE /api/v1/users/skills/wanted/:id
	}
}
