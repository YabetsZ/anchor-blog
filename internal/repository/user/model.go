package userrepo

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"log"
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

// :::::::::   Mapping functions   :::::::::::::::
func ModelToEntity(model *User) entities.User {
	socialLinks := make([]entities.SocialLink, len(model.Profile.SocialLink))

	for index, socialLink := range model.Profile.SocialLink {
		socialLinks[index] = entities.SocialLink{Platform: socialLink.Platform, URL: socialLink.URL}
	}

	posts := make([]string, len(model.UserPosts))
	for index, post := range model.UserPosts {
		posts[index] = post.Hex()
	}

	return entities.User{
		ID:           model.ID.Hex(),
		Username:     model.Username,
		FirstName:    model.FirstName,
		LastName:     model.LastName,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		Role:         model.Role,
		Activated:    model.Activated,
		LastSeen:     model.LastSeen,
		CreatedAt:    model.CreatedAt,
		UpdatedBy:    model.UpdatedBy.Hex(),
		UpdatedAt:    model.UpdatedAt,
		UserPosts:    posts,
		Profile: entities.UserProfile{
			Bio:         model.Profile.Bio,
			PictureURL:  model.Profile.PictureURL,
			SocialLinks: socialLinks,
		},
	}
}

func EntityToModel(ue *entities.User) (*User, error) {
	id, err := primitive.ObjectIDFromHex(ue.ID)
	if err != nil {
		log.Println("invalid user id: ", err.Error())
		return nil, errors.ErrInvalidUserID
	}
	updatedBy, err := primitive.ObjectIDFromHex(ue.UpdatedBy)
	if err != nil {
		log.Println("invalid user id: ", err.Error())
		return nil, errors.ErrInvalidUserID
	}

	userPosts := make([]primitive.ObjectID, len(ue.UserPosts))
	for index, post := range ue.UserPosts {
		userPosts[index], err = primitive.ObjectIDFromHex(post)
		if err != nil {
			log.Println("invalid user id: ", err.Error())
			return nil, errors.ErrInvalidUserID
		}
	}
	socialLinks := make([]SocialLink, len(ue.Profile.SocialLinks))
	for index, socialLink := range ue.Profile.SocialLinks {
		socialLinks[index] = SocialLink{Platform: socialLink.Platform, URL: socialLink.URL}
	}
	return &User{
		ID:           id,
		Username:     ue.Username,
		FirstName:    ue.FirstName,
		LastName:     ue.LastName,
		Email:        ue.Email,
		PasswordHash: ue.PasswordHash,
		Role:         ue.Role,
		Activated:    ue.Activated,
		LastSeen:     ue.LastSeen,
		CreatedAt:    ue.CreatedAt,
		UpdatedAt:    ue.UpdatedAt,
		UpdatedBy:    updatedBy,
		UserPosts:    userPosts,
		Profile: UserProfile{
			Bio:        ue.Profile.Bio,
			PictureURL: ue.Profile.PictureURL,
			SocialLink: socialLinks,
		},
	}, nil
}
