package middleware

import (
	"anchor-blog/api/handler"
	"anchor-blog/pkg/jwtutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or malformed token"})
			return
		}

		tokenString := authHeader[7:]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}
		claim, err := jwtutil.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			handler.HandleHttpError(c, err)
			c.Abort()
			return
		}

		c.Set("UserID", claim.UserID)
		c.Set("Username", claim.Username)
		c.Set("Role", claim.Role)

		c.Next()
	}
}
