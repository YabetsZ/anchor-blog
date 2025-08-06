package handler

import (
	usersvc "anchor-blog/internal/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PasswordResetHandler struct {
	passwordResetService *usersvc.PasswordResetService
}

// NewPasswordResetHandler creates a new password reset handler
func NewPasswordResetHandler(passwordResetService *usersvc.PasswordResetService) *PasswordResetHandler {
	return &PasswordResetHandler{
		passwordResetService: passwordResetService,
	}
}

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the request body for reset password
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ForgotPassword handles POST /api/v1/users/forgot-password
func (h *PasswordResetHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Process forgot password request
	err := h.passwordResetService.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to process forgot password request",
			"details": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset email sent successfully",
		"email":   req.Email,
	})
}

// ResetPassword handles POST /api/v1/users/reset-password
func (h *PasswordResetHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	// Bind and validate request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// Process password reset
	user, err := h.passwordResetService.ResetPassword(c.Request.Context(), req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to reset password",
			"details": err.Error(),
		})
		return
	}

	// Return success response (exclude password hash)
	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}