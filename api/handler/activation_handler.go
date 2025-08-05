package handler

import (
	usersvc "anchor-blog/internal/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ActivationHandler struct {
	activationService *usersvc.ActivationService
}

// NewActivationHandler creates a new activation handler
func NewActivationHandler() *ActivationHandler {
	return &ActivationHandler{
		activationService: usersvc.NewActivationService(),
	}
}

// ActivateAccount handles GET /api/v1/users/activate
func (h *ActivationHandler) ActivateAccount(c *gin.Context) {
	// Extract token from query parameter
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Activation token is required",
		})
		return
	}

	// Verify the activation token
	user, err := h.activationService.VerifyActivation(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid or expired activation token",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Account activated successfully",
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"activated": user.Activated,
			"role":      user.Role,
		},
	})
}
