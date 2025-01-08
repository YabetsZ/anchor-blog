package middleware

import (
	"anchor-blog/internal/errors"
	"anchor-blog/pkg/jwtutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and attaches user info to context
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if it starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must start with Bearer",
			})
			c.Abort()
			return
		}

		// Extract the token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
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

		// Continue to next handler
		c.Next()
	}
}