package handler

import (
	AppError "anchor-blog/internal/errors"
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
	case errors.Is(err, AppError.ErrNotFound),
		errors.Is(err, AppError.ErrUserNotFound):

		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, AppError.ErrInvalidUserID),
		errors.Is(err, AppError.ErrInvalidPostID),
		errors.Is(err, AppError.ErrValidationFailed),
		errors.Is(err, AppError.ErrInvalidToken):

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, AppError.ErrEmailAlreadyExists),
		errors.Is(err, AppError.ErrUsernameTaken):

		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, AppError.ErrInvalidCredentials),
		errors.Is(err, AppError.ErrUnauthorized):

		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	case errors.Is(err, AppError.ErrForbidden),
		errors.Is(err, AppError.ErrUserIsUnverified),
		errors.Is(err, AppError.ErrUserAlreadyAdmin):

		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

	case errors.Is(err, AppError.ErrInternalServer),
		errors.Is(err, AppError.ErrFailedToParse),
		errors.Is(err, AppError.ErrNameCannotEmpty):

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
