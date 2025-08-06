package usersvc

import (
	"anchor-blog/internal/domain/entities"
	"time"
)

type SocialLinkDTO struct {
	Platform string `json:"platform"`
	URL      string `json:"url"`
}

type UserProfileDTO struct {
	Bio         string          `json:"bio"`
	PictureURL  string          `json:"picture_url"`
	SocialLinks []SocialLinkDTO `json:"social_links"`
}

type UserDTO struct {
	ID        string         `json:"id"`
	Username  string         `json:"username"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	Role      string         `json:"role"`
	Password  string         `json:"password"` // this doesn't cause problem! let me know of cases it might cause one
	Activated bool           `json:"activated"`
	LastSeen  time.Time      `json:"last_seen"`
	Profile   UserProfileDTO `json:"profile"`
	UpdatedBy string         `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UserPosts []string       `json:"user_posts"`
}

// :::::::::  Mapping functions  ::::::::::
func EntityToDTO(ue entities.User) UserDTO {
	socialLinks := make([]SocialLinkDTO, len(ue.Profile.SocialLinks))

	for index, socialLink := range ue.Profile.SocialLinks {
		socialLinks[index] = SocialLinkDTO{Platform: socialLink.Platform, URL: socialLink.URL}
	}

	return UserDTO{
		ID:        ue.ID,
		Username:  ue.Username,
		FirstName: ue.FirstName,
		LastName:  ue.LastName,
		Email:     ue.Email,
		Role:      ue.Role,
		Activated: ue.Activated,
		LastSeen:  ue.LastSeen,
		CreatedAt: ue.CreatedAt,
		UpdatedBy: ue.UpdatedBy,
		UpdatedAt: ue.UpdatedAt,
		UserPosts: ue.UserPosts,
		Profile: UserProfileDTO{
			Bio:         ue.Profile.Bio,
			PictureURL:  ue.Profile.PictureURL,
			SocialLinks: socialLinks,
		},
	}
}

func DTOToEntity(dto UserDTO) entities.User {
	socialLinks := make([]entities.SocialLink, len(dto.Profile.SocialLinks))

	for index, socialLink := range dto.Profile.SocialLinks {
		socialLinks[index] = entities.SocialLink{
			Platform: socialLink.Platform,
			URL:      socialLink.URL,
		}
	}

	return entities.User{
		ID:           dto.ID,
		Username:     dto.Username,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Email:        dto.Email,
		Role:         dto.Role,
		PasswordHash: dto.Password,
		Activated:    dto.Activated,
		LastSeen:     dto.LastSeen,
		CreatedAt:    dto.CreatedAt,
		UpdatedBy:    dto.UpdatedBy,
		UpdatedAt:    dto.UpdatedAt,
		UserPosts:    dto.UserPosts,
		Profile: entities.UserProfile{
			Bio:         dto.Profile.Bio,
			PictureURL:  dto.Profile.PictureURL,
			SocialLinks: socialLinks,
		},
	}
}
