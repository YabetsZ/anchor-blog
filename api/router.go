package api

import (
	"anchor-blog/api/handler"
	"anchor-blog/api/handler/user"
	"anchor-blog/internal/repository/gemini"
	contentsvc "anchor-blog/internal/service/content"

	"anchor-blog/api/handler/content"

	"anchor-blog/api/middleware"

	"anchor-blog/config"

	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, userHandler *user.UserHandler, postHandler *handler.PostHandler) *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	v1 := router.Group("/api/v1")

	// Public routes
	public := v1.Group("")
	{
		// Auth routes
		public.POST("/user/register", userHandler.Register) // ✔️
		public.POST("/user/login", userHandler.Login)       // ✔️
		public.POST("/refresh", userHandler.Refresh)        // ✔️

		// Post routes
		public.GET("/posts/:id", postHandler.GetByID)
		public.GET("/posts", postHandler.List)

		// Account activation
		// public.GET("activate", activationHandler.ActivateAccount)
	}

	private := v1.Group("")
	private.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))
	{
		private.POST("/posts", postHandler.Create)
	}

	contentRepo := gemini.NewGeminiRepo(cfg.GenAI.GeminiAPIKey, cfg.GenAI.GeminiModel)
	contentUsecase := contentsvc.NewContentUsecase(contentRepo)
	contentHandler := content.NewContentHandler(contentUsecase)

	aiGenerate := router.Group("/api/v1/ai")
	aiGenerate.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))

	{
		aiGenerate.POST("/generate", contentHandler.GenerateContent)
	}

	return router
}
