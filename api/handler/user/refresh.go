package user

import (
	"anchor-blog/api/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uh *UserHandler) Refresh(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
		return
	}

	tokenString := authHeader[7:]

	loginResponse, err := uh.UserService.Refresh(tokenString)
	if err != nil {
		handler.HandleHttpError(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}
