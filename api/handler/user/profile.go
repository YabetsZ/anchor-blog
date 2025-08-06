package user

import (
	"anchor-blog/api/handler"
	usersvc "anchor-blog/internal/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	profileService *usersvc.ProfileService
}

func NewProfileHandler(profileService *usersvc.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the profile information of the authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} usersvc.ProfileResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 404 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/v1/profile [get]
func (ph *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.HandleError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	profile, err := ph.profileService.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		handler.HandleError(c, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update the profile information of the authenticated user
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body usersvc.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} usersvc.ProfileResponse
// @Failure 400 {object} handler.ErrorResponse
// @Failure 401 {object} handler.ErrorResponse
// @Failure 500 {object} handler.ErrorResponse
// @Router /api/v1/profile [put]
func (ph *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.HandleError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req usersvc.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handler.HandleError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	profile, err := ph.profileService.UpdateUserProfile(c.Request.Context(), userID.(string), &req)
	if err != nil {
		handler.HandleError(c, http.StatusInternalServerError, "Failed to update profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
		"message": "Profile updated successfully",
	})
}