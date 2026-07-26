package routes

import (
	"github.com/Samueljr-web/klinme-api/handlers"
	"github.com/Samueljr-web/klinme-api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	api := r.Group("/api")
	{
		// health check
		api.GET("/ping", handlers.Ping)
		api.POST("/webhooks/clerk", handlers.ClerkWebhook)

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			files := protected.Group("/files")
			{
				files.POST("/upload", handlers.UploadFile)
			}

			users := protected.Group("/users")
			{
				// users.POST("/", handlers.CreateUser)
				users.GET("/:id", handlers.GetUser)
			}
		}

	}
}
