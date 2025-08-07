package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole creates a middleware that checks if the user has the required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Check if user has the required role
		if userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
				"required_role": requiredRole,
				"user_role": userRole,
			})
			c.Abort()
			return
		}

		// User has the required role, continue
		c.Next()
	}
}

// RequireAdmin is a convenience middleware for admin-only routes
func RequireAdmin() gin.HandlerFunc {
	return RequireRole("admin")
}

// RequireUser is a convenience middleware for user-level routes (excludes unverified)
func RequireUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Allow both "user" and "admin" roles, but not "unverified"
		if userRole != "user" && userRole != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Account not verified or insufficient permissions",
				"user_role": userRole,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}