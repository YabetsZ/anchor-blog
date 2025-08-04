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
