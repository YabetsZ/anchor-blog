package api

import (
	"anchor-blog/api/handler"
  "anchor-blog/api/handler/user"

	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
	// g := router.Group("")
	// userHandler := user.NewUserHandler(usersvc.NewUserServices()) // insert it after segni has done the job
	// UserRoutes(g, userHandler)

	// Initialize handlers
	activationHandler := handler.NewActivationHandler()

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User activation route
		v1.GET("/users/activate", activationHandler.ActivateAccount)
	}

	return router
}

func UserRoutes(rg *gin.RouterGroup, handler *user.UserHandler) {
	// Public routes
	rg.POST("/login", handler.Login)
	rg.POST("/refresh", handler.Refresh)

}
