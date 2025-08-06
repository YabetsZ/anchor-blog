package usersvc

import (
	"anchor-blog/internal/domain/entities"
	errorr "anchor-blog/internal/errors"
	"context"
	"time"
)

type ProfileService struct {
	userRepo entities.IUserRepository
}

func NewProfileService(userRepo entities.IUserRepository) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
	}
}

type UpdateProfileRequest struct {
	Bio         *string          `json:"bio,omitempty"`
	PictureURL  *string          `json:"picture_url,omitempty"`
	SocialLinks *[]SocialLinkDTO `json:"social_links,omitempty"`
}

type ProfileResponse struct {
	Bio         string          `json:"bio"`
	PictureURL  string          `json:"picture_url"`
	SocialLinks []SocialLinkDTO `json:"social_links"`
}

func (ps *ProfileService) GetUserProfile(ctx context.Context, userID string) (*ProfileResponse, error) {
	user, err := ps.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	socialLinks := make([]SocialLinkDTO, len(user.Profile.SocialLinks))
	for i, link := range user.Profile.SocialLinks {
		socialLinks[i] = SocialLinkDTO{
			Platform: link.Platform,
			URL:      link.URL,
		}
	}

	return &ProfileResponse{
		Bio:         user.Profile.Bio,
		PictureURL:  user.Profile.PictureURL,
		SocialLinks: socialLinks,
	}, nil
}

func (ps *ProfileService) UpdateUserProfile(ctx context.Context, userID string, req *UpdateProfileRequest) (*ProfileResponse, error) {
	// Get current user to preserve existing data
	currentUser, err := ps.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Create update entity with only the fields that should be updated
	updateUser := &entities.User{
		ID: userID,
		Profile: entities.UserProfile{
			Bio:        currentUser.Profile.Bio,
			PictureURL: currentUser.Profile.PictureURL,
			SocialLinks: currentUser.Profile.SocialLinks,
		},
		UpdatedAt: time.Now(),
	}

	// Update only provided fields
	if req.Bio != nil {
		updateUser.Profile.Bio = *req.Bio
	}

	if req.PictureURL != nil {
		updateUser.Profile.PictureURL = *req.PictureURL
	}

	if req.SocialLinks != nil {
		socialLinks := make([]entities.SocialLink, len(*req.SocialLinks))
		for i, link := range *req.SocialLinks {
			socialLinks[i] = entities.SocialLink{
				Platform: link.Platform,
				URL:      link.URL,
			}
		}
		updateUser.Profile.SocialLinks = socialLinks
	}

	// Update in repository
	err = ps.userRepo.EditUserByID(ctx, userID, updateUser)
	if err != nil {
		return nil, errorr.ErrInternalServer
	}

	// Return updated profile
	return ps.GetUserProfile(ctx, userID)
}