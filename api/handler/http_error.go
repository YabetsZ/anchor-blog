package handler

import (
	errs "anchor-blog/internal/errors"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Success bool   `json:"success"`
}

func HandleHttpError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errs.ErrNotFound),
		errors.Is(err, errs.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, errs.ErrInvalidUserID),
		errors.Is(err, errs.ErrInvalidPostID),
		errors.Is(err, errs.ErrValidationFailed),
		errors.Is(err, errs.ErrInvalidToken):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, errs.ErrEmailAlreadyExists),
		errors.Is(err, errs.ErrUsernameTaken):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, errs.ErrInvalidCredentials),
		errors.Is(err, errs.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	case errors.Is(err, errs.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, errs.ErrInternalServer):
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	case errors.Is(err, errs.ErrNameCannotEmpty):
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		log.Printf("An unexpected error occurred: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error occurred"})
	}
}

func HandleError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, ErrorResponse{
		Error:   message,
		Success: false,
	})
}
