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

	userGroup := router.Group("/api/v1/user")

	UserRoutes(cfg, userGroup, userHandler)

	v1Auth := router.Group("/api/v1")
	v1Auth.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))
	{
		// ... your existing authenticated routes for profile, etc.

		// Post routes
		v1Auth.POST("/posts", postHandler.Create)
	}

	// Initialize handlers
	activationHandler := handler.NewActivationHandler()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User activation route
		v1.GET("activate", activationHandler.ActivateAccount)
	}

	// Public routes that don't need auth
	v1Public := router.Group("/api/v1")
	{
		v1Public.GET("/posts/:id", postHandler.GetByID)
		v1Public.GET("/posts", postHandler.List)
	}

	return router
}

func UserRoutes(cfg *config.Config, rg *gin.RouterGroup, handler *user.UserHandler) {
	// Public routes

	public := rg.Group("")
	public.POST("/login", handler.Login)
	public.POST("/refresh", handler.Refresh)

	private := rg.Group("")
	private.Use(middleware.AuthMiddleware(cfg.JWT.AccessTokenSecret))

}

func PostRoutes() {

	rg.POST("/login", handler.Login)
	rg.POST("/refresh", handler.Refresh)
	rg.POST("/register", handler.Register)


}
