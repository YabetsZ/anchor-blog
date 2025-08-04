package api

import (
	"anchor-blog/api/handler/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})
	// g := router.Group("")
	// userHandler := user.NewUserHandler(usersvc.NewUserServices()) // insert it after segni has done the job
	// UserRoutes(g, userHandler)

	return router
}

func UserRoutes(rg *gin.RouterGroup, handler *user.UserHandler) {
	// Public routes
	rg.POST("/login", handler.Login)
	rg.POST("/refresh", handler.Refresh)

}
