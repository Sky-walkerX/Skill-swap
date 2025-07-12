package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/repository"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes by delegating to specific route files
func SetupRoutes(api *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, *cfg)

	// Setup route groups
	SetupAuthRoutes(api, authService, cfg)
	SetupUserRoutes(api, userService, cfg)
	SetupSkillRoutes(api, cfg)
	SetupSwapRoutes(api, cfg)
	SetupAdminRoutes(api, cfg)
}
