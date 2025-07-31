package Domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialLink struct {
	Platform string `bson:"platform" json:"platform"`
	URL      string `bson:"url" json:"url"`
}

type UserProfile struct {
	Bio         string       `bson:"bio" json:"bio"`
	PictureURL  string       `bson:"picture_url" json:"picture_url"`
	SocialLinks []SocialLink `bson:"social_links" json:"social_links"`
}

type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Username     string               `bson:"username" json:"username"`
	FirstName    string               `bson:"first_name" json:"first_name"`
	LastName     string               `bson:"last_name" json:"last_name"`
	Email        string               `bson:"email" json:"email"`
	PasswordHash string               `bson:"password_hash" json:"-"`
	Role         string               `bson:"role" json:"role"` // "user", "admin", "unverified"
	Activated    bool                 `bson:"activated" json:"activated"`
	LastSeen     time.Time            `bson:"last_seen" json:"last_seen"`
	Profile      UserProfile          `bson:"profile" json:"profile"`
	UpdatedBy    primitive.ObjectID   `bson:"updated_by" json:"updated_by"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
	UserPosts    []primitive.ObjectID `bson:"user_posts" json:"user_posts"`
}
