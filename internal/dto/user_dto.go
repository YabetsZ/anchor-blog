package dto

import "time"

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
	Activated bool           `json:"activated"`
	LastSeen  time.Time      `json:"last_seen"`
	Profile   UserProfileDTO `json:"profile"`
	UpdatedBy string         `json:"updated_by"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UserPosts []string       `json:"user_posts"`
}
