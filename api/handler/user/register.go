package user

import (
	"anchor-blog/api/handler"
	usersvc "anchor-blog/internal/service/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (uh *UserHandler) Register(c *gin.Context) {
	var input usersvc.UserDTO
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	response, err := uh.UserService.Register(c, &input)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
