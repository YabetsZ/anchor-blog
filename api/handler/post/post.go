package post

import (
	"anchor-blog/api/handler"
	"anchor-blog/internal/service"
	viewsvc "anchor-blog/internal/service/view"
	"anchor-blog/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService        *service.PostService
	viewTrackingService *viewsvc.ViewTrackingService
}

func NewPostHandler(ps *service.PostService, vts *viewsvc.ViewTrackingService) *PostHandler {
	return &PostHandler{
		postService:        ps,
		viewTrackingService: vts,
	}
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
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	postID := c.Param("id")

	// Track the view with IP-based throttling
	if h.viewTrackingService != nil {
		clientIP := utils.GetClientIP(c)
		err := h.viewTrackingService.TrackView(c.Request.Context(), postID, clientIP)
		if err != nil {
			// Log the error but don't fail the request
			// View tracking is not critical for post retrieval
			c.Header("X-View-Tracking-Error", "true")
		}
	}

	post, err := h.postService.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		handler.HandleHttpError(c, err)
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
		handler.HandleHttpError(c, err)
		return
	}
	res := make([]*PostDTO, len(posts))
	for idx, post := range posts {
		res[idx] = MapPostToDTO(post)
	}

	c.JSON(http.StatusOK, posts)
}

// GetPopularPosts returns posts ordered by view count
func (h *PostHandler) GetPopularPosts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	if limit <= 0 || limit > 100 {
		limit = 10
	}

	posts, err := h.viewTrackingService.GetPopularPosts(c.Request.Context(), limit)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	res := make([]*PostDTO, len(posts))
	for idx, post := range posts {
		res[idx] = MapPostToDTO(post)
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": res,
		"count": len(res),
	})
}

// GetViewStats returns view statistics
func (h *PostHandler) GetViewStats(c *gin.Context) {
	totalViews, err := h.viewTrackingService.GetTotalViews(c.Request.Context())
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_views": totalViews,
	})
}

// GetPostViewCount returns the view count for a specific post
func (h *PostHandler) GetPostViewCount(c *gin.Context) {
	postID := c.Param("id")

	viewCount, err := h.viewTrackingService.GetViewCount(c.Request.Context(), postID)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id":    postID,
		"view_count": viewCount,
	})
}
