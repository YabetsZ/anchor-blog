package viewsvc

import (
	"context"
	"fmt"
	"log"
	"time"

	"anchor-blog/internal/domain/entities"
	redisclient "anchor-blog/pkg/redis"
)

type ViewTrackingService struct {
	redisClient     *redisclient.Client
	postRepo        entities.IPostRepository
	viewTrackingTTL time.Duration
}

func NewViewTrackingService(redisClient *redisclient.Client, postRepo entities.IPostRepository, ttlSeconds int) *ViewTrackingService {
	return &ViewTrackingService{
		redisClient:     redisClient,
		postRepo:        postRepo,
		viewTrackingTTL: time.Duration(ttlSeconds) * time.Second,
	}
}

// TrackView handles view tracking with IP-based throttling
func (vts *ViewTrackingService) TrackView(ctx context.Context, postID, ipAddress string) error {
	// Create Redis key for this IP-Post combination
	viewKey := fmt.Sprintf("post_view:%s:%s", postID, ipAddress)
	
	// Check if this IP has already viewed this post within the TTL period
	exists, err := vts.redisClient.Exists(ctx, viewKey)
	if err != nil {
		log.Printf("Error checking Redis key existence: %v", err)
		// Continue with view tracking even if Redis fails
	}

	// If the key doesn't exist, this is a new view
	if !exists {
		// Set the Redis key with TTL to prevent duplicate views
		err = vts.redisClient.SetWithExpiration(ctx, viewKey, "viewed", vts.viewTrackingTTL)
		if err != nil {
			log.Printf("Error setting Redis key: %v", err)
			// Continue with view tracking even if Redis fails
		}

		// Increment the view count in the database
		err = vts.postRepo.IncrementViewCount(ctx, postID)
		if err != nil {
			log.Printf("Error incrementing view count in database: %v", err)
			return err
		}

		log.Printf("View tracked for post %s from IP %s", postID, ipAddress)
	} else {
		log.Printf("Duplicate view prevented for post %s from IP %s", postID, ipAddress)
	}

	return nil
}

// GetViewCount retrieves the current view count for a post
func (vts *ViewTrackingService) GetViewCount(ctx context.Context, postID string) (int, error) {
	return vts.postRepo.GetViewCount(ctx, postID)
}

// GetTotalViews gets total views across all posts
func (vts *ViewTrackingService) GetTotalViews(ctx context.Context) (int64, error) {
	return vts.postRepo.GetTotalViews(ctx)
}

// GetPopularPosts gets posts ordered by view count
func (vts *ViewTrackingService) GetPopularPosts(ctx context.Context, limit int) ([]*entities.Post, error) {
	return vts.postRepo.GetPostsByViewCount(ctx, limit)
}

// ResetViewTracking removes all view tracking data (admin function)
func (vts *ViewTrackingService) ResetViewTracking(ctx context.Context, postID string) error {
	// This would require scanning Redis keys, which is expensive
	// For now, we'll just reset the database count
	return vts.postRepo.ResetViewCount(ctx, postID)
}