package user

import (
	"anchor-blog/api/handler"
	usersvc "anchor-blog/internal/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *usersvc.UserServices
}

func NewUserHandler(us *usersvc.UserServices) *UserHandler {
	return &UserHandler{UserService: us}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (uh *UserHandler) Login(c *gin.Context) {
	var input LoginRequest
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	response, err := uh.UserService.Login(input.Username, input.Password)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

// Profile methods
func (uh *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.HandleError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	profile, err := uh.UserService.ProfileService.GetUserProfile(c.Request.Context(), userID.(string))
	if err != nil {
		handler.HandleError(c, http.StatusInternalServerError, "Failed to get profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

func (uh *UserHandler) UpdateProfile(c *gin.Context) {
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

	profile, err := uh.UserService.ProfileService.UpdateUserProfile(c.Request.Context(), userID.(string), &req)
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
