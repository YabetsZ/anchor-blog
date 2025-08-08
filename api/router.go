package api

import (
	"anchor-blog/api/handler"
	"anchor-blog/api/handler/content"
	"anchor-blog/api/handler/post"
	"anchor-blog/api/handler/user"
	"anchor-blog/api/middleware"
	"anchor-blog/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config, userHandler *user.UserHandler, postHandler *post.PostHandler, activationHandler *handler.ActivationHandler, passwordResetHandler *handler.PasswordResetHandler, contentHandler *content.ContentHandler) *gin.Engine {
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

		// User activation and password reset routes
		public.GET("/users/activate", activationHandler.ActivateAccount)
		public.POST("/users/forgot-password", passwordResetHandler.ForgotPassword)
		public.POST("/users/reset-password", passwordResetHandler.ResetPassword)

		// Post routes
		public.GET("/posts/:id", postHandler.GetByID)                // ✔️
		public.GET("/posts", postHandler.List)                       // ✔️
		public.GET("/posts/popular", postHandler.GetPopularPosts)    // ✔️
		public.GET("/posts/search", postHandler.SearchPosts)         // ✔️
		public.GET("/posts/filter", postHandler.FilterPosts)         // ✔️
		public.GET("/posts/:id/views", postHandler.GetPostViewCount) // ✔️
		public.GET("/stats/views", postHandler.GetViewStats)         // ✔️
	}

	private := v1.Group("")
	private.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))
	{
		// Post routes
		private.POST("/posts", postHandler.Create)           // ✔️
		private.PUT("/posts/:id", postHandler.UpdatePost)    // ✔️
		private.DELETE("/posts/:id", postHandler.DeletePost) // ✔️

		// Post interaction routes
		private.POST("/posts/:id/like", postHandler.LikePost)                // ✔️
		private.DELETE("/posts/:id/like", postHandler.UnlikePost)            // ✔️
		private.POST("/posts/:id/dislike", postHandler.DislikePost)          // ✔️
		private.DELETE("/posts/:id/dislike", postHandler.UndislikePost)      // ✔️
		private.GET("/posts/:id/like-status", postHandler.GetPostLikeStatus) // ✔️

		// Profile routes
		private.GET("/user/profile", userHandler.GetProfile)
		private.PUT("/user/profile", userHandler.UpdateProfile)

		// Admin routes
		private.PATCH("/admin/users/:id/promote", middleware.RequireSuperadmin(), userHandler.PromoteUser) // consider using superadmin
		private.PATCH("/admin/users/:id/demote", middleware.RequireSuperadmin(), userHandler.DemoteUser)   // consider using superadmin

		// Auth routes
		private.POST("/logout", userHandler.Logout) // ✔️
	}

	// AI Content Generation routes
	aiGenerate := router.Group("/api/v1/ai")
	aiGenerate.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))
	{
		aiGenerate.POST("/generate", contentHandler.GenerateContent)
	}

	return router
}
