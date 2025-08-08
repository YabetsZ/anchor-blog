package user

import (
	"anchor-blog/api/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) PromoteUser(c *gin.Context) {
	targetUserID := c.Param("id")
	adminID := c.GetString("user_id")

	err := h.UserService.PromoteUserToAdmin(c.Request.Context(), adminID, targetUserID)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User successfully promoted to admin"})
}

func (h *UserHandler) DemoteUser(c *gin.Context) {
	targetUserID := c.Param("id")
	adminID := c.GetString("user_id")

	err := h.UserService.DemoteAdminToUser(c.Request.Context(), adminID, targetUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Admin successfully demoted to user"})
}
