package middleware

import (

	"anchor-blog/internal/errors"

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


		// Validate the token using JWT utilities
		claims, err := jwtutil.ValidateToken(token, jwtSecret)
		if err != nil {
			if err == errors.ErrInvalidToken {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid or expired token",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token validation failed",
				})
			}
			c.Abort()
			return
		}


		// Extract user info from JWT claims and attach to context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
