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
	id, err := uh.UserService.Register(c.Request.Context(), &input)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}
	response := registerResponse{ID: id}
	c.JSON(http.StatusOK, response)
}
