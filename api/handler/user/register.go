package user

import (
	"anchor-blog/api/handler"
	usersvc "anchor-blog/internal/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerResponse struct {
	ID string `json:"id"`
}

func (uh *UserHandler) Register(c *gin.Context) {
	var input usersvc.UserDTO
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	
	// Create the user account
	id, err := uh.UserService.Register(c.Request.Context(), &input)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	// Get the created user to generate activation token
	user, err := uh.UserService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	// Generate and send activation token (only for non-superadmin users)
	if user.Role != "superadmin" && uh.ActivationService != nil {
		err = uh.ActivationService.SendActivationEmail(c.Request.Context(), user)
		if err != nil {
			// Log the error but don't fail the registration
			// The user is created, they just won't get an activation email
			c.Header("X-Activation-Warning", "Failed to send activation email")
		}
	}

	response := registerResponse{ID: id}
	c.JSON(http.StatusOK, response)
}
