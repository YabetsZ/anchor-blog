package content

import (
	"anchor-blog/internal/domain/entities"
	AppError "anchor-blog/internal/errors"
	contentsvc "anchor-blog/internal/service/content"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	uc contentsvc.ContentUsecase
}

func NewContentHandler(uc contentsvc.ContentUsecase) *ContentHandler {
	return &ContentHandler{uc: uc}
}

func (h *ContentHandler) GenerateContent(c *gin.Context) {
	var req entities.ContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "bad request",
		})
		return
	}

	if req.WordCount < 10 || req.WordCount > 200 {
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Word count must be between 10 and 200",
		})
		return
	}

	resp, err := h.uc.GenerateContent(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Add security headers
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-Frame-Options", "DENY")

	c.JSON(http.StatusOK, resp)
}

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, AppError.ErrInvalidInput):
		c.JSON(http.StatusBadRequest, entities.ErrorResponse{
			Error: "Invalid input parameters",
		})
	case errors.Is(err, AppError.ErrContentBlocked),
		errors.Is(err, AppError.ErrIllegalContent):
		c.JSON(http.StatusUnprocessableEntity, entities.ErrorResponse{
			Error: "Content blocked by safety filters",
		})
	case errors.Is(err, context.DeadlineExceeded):
		c.JSON(http.StatusGatewayTimeout, entities.ErrorResponse{
			Error: "Processing timeout",
		})
	default:
		c.JSON(http.StatusInternalServerError, entities.ErrorResponse{
			Error: "Internal server error",
		})
	}
}
