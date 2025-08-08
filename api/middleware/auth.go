package middleware

import (
	"anchor-blog/internal/errors"
	"log"
	"strings"

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
		log.Println(authHeader)
		tokenString := authHeader[7:]
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			c.Abort()
			return
		}
		log.Println("Number of toke segment", len(strings.Split(tokenString, ".")))
		log.Println(strings.Split(tokenString, "."))

		// Validate the token using JWT utilities
		claims, err := jwtutil.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			if err == errors.ErrInvalidToken {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Invalid or expired token",
				})
			} else {
				log.Println(err.Error())
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Token validation failed",
				})
			}
			c.Abort()
			return
		}

		// Extract user info from JWT claims and attach to context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}
