package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SocialLink struct {
	Platform string `bson:"platform"`
	URL      string `bson:"url"`
}

type UserProfile struct {
	Bio        string       `bson:"bio"`
	PictureURL string       `bson:"picture_url"`
	SocialLink []SocialLink `bson:"social_links"`
}

type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty"`
	Username     string               `bson:"username"`
	FirstName    string               `bson:"first_name"`
	LastName     string               `bson:"last_name"`
	Email        string               `bson:"email"`
	PasswordHash string               `bson:"password_hash"`
	Role         string               `bson:"role"`
	Activated    bool                 `bson:"activated"`
	LastSeen     time.Time            `bson:"last_seen"`
	Profile      UserProfile          `bson:"profile"`
	UpdatedBy    primitive.ObjectID   `bson:"updated_by"`
	CreatedAt    time.Time            `bson:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at"`
	UserPosts    []primitive.ObjectID `bson:"user_posts"`
}
