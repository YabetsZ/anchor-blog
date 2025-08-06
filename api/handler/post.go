package handler

import (
	"anchor-blog/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(ps *service.PostService) *PostHandler {
	return &PostHandler{postService: ps}
}

type CreatePostRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Tags    []string `json:"tags"`
}

func (h *PostHandler) Create(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get author ID from context, set by the AuthMiddleware
	authorIDHex, exists := c.Get("UserID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	post, err := h.postService.CreatePost(c.Request.Context(), req.Title, req.Content, authorIDHex.(string), req.Tags)
	if err != nil {
		HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	postID := c.Param("id")

	post, err := h.postService.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.ParseInt(pageStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	posts, err := h.postService.ListPosts(c.Request.Context(), page, limit)
	if err != nil {
		HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}
