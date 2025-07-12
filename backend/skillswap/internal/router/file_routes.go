package router

import (
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/app/service"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/config"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/file"
	"github.com/Sky-walkerX/Skill-swap/backend/skillswap/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupFileRoutes(api *gin.RouterGroup, fileUploadService *service.FileUploadService, cfg *config.Config) {
	fileHandler := file.NewHandler(fileUploadService)

	// File routes
	files := api.Group("/files")
	{
		// Public file access (for serving uploaded photos)
		files.GET("/users/:user_id/photo", fileHandler.GetUserPhoto) // GET /api/files/users/:user_id/photo

		// Protected file routes
		protected := files.Group("/users")
		protected.Use(middleware.JWTAuth(*cfg))
		{
			protected.POST("/photo", fileHandler.UploadUserPhoto)         // POST /api/files/users/photo
			protected.DELETE("/photo", fileHandler.DeleteUserPhoto)       // DELETE /api/files/users/photo
			protected.GET("/:user_id/info", fileHandler.GetUserPhotoInfo) // GET /api/files/users/:user_id/info
		}
	}
}
