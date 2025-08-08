package post

import (
	"anchor-blog/api/handler"
	"anchor-blog/internal/domain/entities"
	postsvc "anchor-blog/internal/service/post"
	viewsvc "anchor-blog/internal/service/view"
	"anchor-blog/pkg/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService         *postsvc.PostService
	viewTrackingService *viewsvc.ViewTrackingService
}

func NewPostHandler(ps *postsvc.PostService, vts *viewsvc.ViewTrackingService) *PostHandler {
	return &PostHandler{
		postService:         ps,
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
	authorIDHex, exists := c.Get("user_id")
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

// UpdatePost updates an existing post
func (h *PostHandler) UpdatePost(c *gin.Context) {
	postID := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user owns the post (basic authorization)
	existingPost, err := h.postService.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	if existingPost.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		return
	}

	updatedPost, err := h.postService.UpdatePost(c.Request.Context(), postID, req.Title, req.Content, req.Tags)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, MapPostToDTO(updatedPost))
}

// DeletePost deletes a post
func (h *PostHandler) DeletePost(c *gin.Context) {
	postID := c.Param("id")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Check if user owns the post (basic authorization)
	existingPost, err := h.postService.GetPostByID(c.Request.Context(), postID)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	if existingPost.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}

	err = h.postService.DeletePost(c.Request.Context(), postID)
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// SearchPosts searches for posts
func (h *PostHandler) SearchPosts(c *gin.Context) {
	query := c.Query("q")
	searchType := c.DefaultQuery("type", "title") // title or author
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	page, _ := strconv.ParseInt(pageStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	posts, err := h.postService.SearchPosts(c.Request.Context(), query, searchType, page, limit)
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
		"query": query,
		"type":  searchType,
	})
}

// FilterPosts filters posts by tags or date
func (h *PostHandler) FilterPosts(c *gin.Context) {
	tagsParam := c.Query("tags")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.ParseInt(pageStr, 10, 64)
	limit, _ := strconv.ParseInt(limitStr, 10, 64)

	var posts []*entities.Post
	var err error

	if tagsParam != "" {
		// Filter by tags
		tags := strings.Split(tagsParam, ",")
		// Trim whitespace from tags
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		posts, err = h.postService.FilterPostsByTags(c.Request.Context(), tags, page, limit)
	} else if startDate != "" && endDate != "" {
		// Filter by date range
		posts, err = h.postService.FilterPostsByDateRange(c.Request.Context(), startDate, endDate, page, limit)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either tags or date range (start_date and end_date) must be provided"})
		return
	}

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

// LikePost likes a post
func (h *PostHandler) LikePost(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.postService.LikePost(c.Request.Context(), postID, userID.(string))
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post liked successfully"})
}

// UnlikePost unlikes a post
func (h *PostHandler) UnlikePost(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.postService.UnlikePost(c.Request.Context(), postID, userID.(string))
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post unliked successfully"})
}

// DislikePost dislikes a post
func (h *PostHandler) DislikePost(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.postService.DislikePost(c.Request.Context(), postID, userID.(string))
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post disliked successfully"})
}

// UndislikePost removes dislike from a post
func (h *PostHandler) UndislikePost(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.postService.UndislikePost(c.Request.Context(), postID, userID.(string))
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post undisliked successfully"})
}

// GetPostLikeStatus gets the like status for a post
func (h *PostHandler) GetPostLikeStatus(c *gin.Context) {
	postID := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	liked, disliked, err := h.postService.GetPostLikeStatus(c.Request.Context(), postID, userID.(string))
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id":  postID,
		"liked":    liked,
		"disliked": disliked,
	})
}
