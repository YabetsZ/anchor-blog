package user

import (
	"anchor-blog/api/handler"
	"anchor-blog/internal/domain/entities"
	AppError "anchor-blog/internal/errors"
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

	loginResponse, err := uh.UserService.Refresh(c.Request.Context(), tokenString)
	if err != nil {
		handler.HandleHttpError(c, err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, loginResponse)
}

func (uh *UserHandler) SetLastSeen(c *gin.Context) {
	userID := c.Param("id")
	err := uh.UserService.SetLastSeen(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{Error: AppError.ErrInternalServer.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
