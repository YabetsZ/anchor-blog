package post

import (
	"anchor-blog/api/handler"
	postsvc "anchor-blog/internal/service/post"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *postsvc.PostService
}

func NewPostHandler(ps *postsvc.PostService) *PostHandler {
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
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetByID(c *gin.Context) {
	postID := c.Param("id")

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

func (h *PostHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.HandleError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	postID := c.Param("id")

	err := h.postService.Delete(c, postID, userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}

func (h *PostHandler) Like(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.HandleError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	postID := c.Param("id")

	liked, err := h.postService.LikePost(c, postID, userID.(string))
	if err != nil {
		log.Println(err.Error())
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *PostHandler) Dislike(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.HandleError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	postID := c.Param("id")

	liked, err := h.postService.DislikePost(c, postID, userID.(string))
	if err != nil {
		handler.HandleHttpError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"disliked": liked})
}

func (h *PostHandler) Edit(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		handler.HandleError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}
	postID := c.Param("id")
	var req PostDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := MapDTOToPost(&req)

	err := h.postService.Update(c, post, postID, userID.(string))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "post edited successfully"})
}
