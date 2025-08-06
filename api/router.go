package api

import (
	"anchor-blog/api/handler"
	"anchor-blog/api/handler/user"
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

	// Initialize activation and password reset handlers
	activationHandler := handler.NewActivationHandler()
	passwordResetHandler := handler.NewPasswordResetHandler()

	// User routes
	userGroup := router.Group("/api/v1/user")
	UserRoutes(cfg, userGroup, userHandler)

	// Authenticated routes
	v1Auth := router.Group("/api/v1")
	v1Auth.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))
	{
		// Post routes
		v1Auth.POST("/posts", postHandler.Create)
	}

	// Public routes that don't need auth
	v1Public := router.Group("/api/v1")
	{
		// Post routes
		v1Public.GET("/posts/:id", postHandler.GetByID)
		v1Public.GET("/posts", postHandler.List)
		
		// User activation and password reset routes (public)
		v1Public.GET("/users/activate", activationHandler.ActivateAccount)
		v1Public.POST("/users/forgot-password", passwordResetHandler.ForgotPassword)
		v1Public.POST("/users/reset-password", passwordResetHandler.ResetPassword)
	}

	return router
}

func UserRoutes(cfg *config.Config, rg *gin.RouterGroup, handler *user.UserHandler) {
	// Public routes
	public := rg.Group("")
	public.POST("/login", handler.Login)
	public.POST("/refresh", handler.Refresh)
	public.POST("/register", handler.Register)

	// Private routes
	private := rg.Group("")
	private.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))
}