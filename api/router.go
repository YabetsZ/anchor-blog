package api

import (
	"anchor-blog/api/handler"
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
