package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/gin-gonic/gin"
)

// SetupSkillRoutes configures all skill-related routes
func SetupSkillRoutes(api *gin.RouterGroup, cfg *config.Config) {
	// TODO: Implement skill routes
	// - GET /skills - List all available skills
	// - POST /skills - Create new skill (admin only)
	// - PUT /skills/:id - Update skill (admin only)
	// - DELETE /skills/:id - Delete skill (admin only)
	// - POST /users/skills/offered - Add skill to user's offered skills
	// - DELETE /users/skills/offered/:id - Remove offered skill
	// - POST /users/skills/wanted - Add skill to user's wanted skills
	// - DELETE /users/skills/wanted/:id - Remove wanted skill

	_ = api // Unused for now
	_ = cfg // Unused for now
}
