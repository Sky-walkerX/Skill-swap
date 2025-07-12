package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/admin"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/repository"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/availability"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/rating"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/skill"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/swap"
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
	skillService := service.NewSkillService(db)
	swapService := service.NewSwapService(db)
	ratingService := service.NewRatingService(db)
	adminService := service.NewAdminService(db)
	availabilityService := service.NewAvailabilityService(db)
	notificationService := service.NewNotificationService(db)
	searchService := service.NewSearchService(db)
	fileUploadService := service.NewFileUploadService(db, *cfg)

	// Initialize handlers
	skillHandler := skill.NewHandler(skillService)
	swapHandler := swap.NewHandler(swapService)
	ratingHandler := rating.NewHandler(ratingService)
	adminHandler := admin.NewHandler(adminService)
	availabilityHandler := availability.NewHandler(availabilityService)

	// Setup route groups
	SetupAuthRoutes(api, authService, cfg)
	SetupUserRoutes(api, userService, cfg)
	SetupSkillRoutes(api, cfg, skillHandler)
	SetupSwapRoutes(api, cfg, swapHandler)
	SetupRatingRoutes(api, cfg, ratingHandler)
	SetupAvailabilityRoutes(api, cfg, availabilityHandler)
	SetupAdminRoutes(api, cfg, skillHandler, adminHandler)
	SetupNotificationRoutes(api, notificationService, cfg)
	SetupSearchRoutes(api, searchService, cfg)
	SetupFileRoutes(api, fileUploadService, cfg)
}
