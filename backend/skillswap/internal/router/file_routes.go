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
		// Public file access (for serving uploaded files)
		files.GET("/users/:user_id/:filename", fileHandler.GetUserPhoto)        // GET /api/files/users/:user_id/:filename
		
		// Protected file routes
		protected := files.Group("/users")
		protected.Use(middleware.JWTAuth(*cfg))
		{
			protected.POST("/photo", fileHandler.UploadUserPhoto)               // POST /api/files/users/photo
			protected.DELETE("/photo", fileHandler.DeleteUserPhoto)            // DELETE /api/files/users/photo
			protected.GET("/:user_id/:filename/info", fileHandler.GetUserPhotoInfo) // GET /api/files/users/:user_id/:filename/info
		}
	}
}
