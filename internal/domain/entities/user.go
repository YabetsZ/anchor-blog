package entities

import (
	"time"
)

const (
	RoleAdmin      = "admin"
	RoleUnverified = "unverified"
	RoleUser       = "user"
	RoleSuperadmin = "superadmin"
)

type SocialLink struct {
	Platform string
	URL      string
}

type UserProfile struct {
	Bio         string
	PictureURL  string
	SocialLinks []SocialLink
}

type User struct {
	ID           string
	Username     string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	Role         string // unverified, user, admin, superadmin
	Activated    bool
	LastSeen     time.Time
	Profile      UserProfile
	UpdatedBy    string // This should be for who changed the role
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserPosts    []string // this will be depricated
}
