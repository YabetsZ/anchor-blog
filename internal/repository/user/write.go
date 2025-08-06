package userrepo

import (
	"anchor-blog/internal/domain/entities"
	"context"
	"log"
	"time"

	errorr "anchor-blog/internal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ur *userRepository) CreateUser(ctx context.Context, user *entities.User) (string, error) {
	userDoc, err := EntityToModel(user)
	userDoc.ID = primitive.NewObjectID()
	userDoc.UpdatedBy = userDoc.ID

	if err != nil {
		log.Printf("error while transfer user entity to user model %v", err.Error())
		return "", errorr.ErrInternalServer
	}
	_, err = ur.collection.InsertOne(ctx, userDoc)
	if err != nil {
		log.Printf("error while create new user %v", err.Error())
		return "", errorr.ErrInternalServer
	}
	return userDoc.ID.Hex(), nil
}
func (ur *userRepository) DeleteUserByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errorr.ErrInternalServer
	}
	filter := bson.M{"id": objID}
	_, err = ur.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("error while delete user %v", err.Error())
		return errorr.ErrInternalServer
	}
	return nil
}

func (ur *userRepository) SetLastSeen(ctx context.Context, id string, timestamp time.Time) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error while cast id to object id %v", err.Error())
		return errorr.ErrInternalServer
	}
	filter := bson.M{"id": objID}

	var foundUser User

	err = ur.collection.FindOne(ctx, filter).Decode(foundUser)
	if err != nil {
		log.Printf("error while find user %v", err.Error())
		return errorr.ErrInternalServer
	}
	foundUser.LastSeen = timestamp
	_, err = ur.collection.UpdateOne(ctx, filter, foundUser)
	if err != nil {
		return errorr.ErrInternalServer
	}
	return err
}

func (ur *userRepository) EditUserByID(ctx context.Context, id string, user *entities.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error while cast id to object id %v", err.Error())
		return errorr.ErrInternalServer
	}

	filter := bson.M{"_id": objID}

	var foundUserModel User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUserModel)
	if err != nil {
		log.Printf("error while find user %v", err.Error())
		return errorr.ErrInternalServer
	}

	foundUser := ModelToEntity(&foundUserModel)

	// Update fields if new data is provided
	if user.Username != "" {
		foundUser.Username = user.Username
	}
	if user.FirstName != "" {
		foundUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		foundUser.LastName = user.LastName
	}

	if !IsUserProfileEmpty(user.Profile) {
		if user.Profile.Bio != "" {
			foundUser.Profile.Bio = user.Profile.Bio
		}
		if user.Profile.PictureURL != "" {
			foundUser.Profile.PictureURL = user.Profile.PictureURL
		}
		if len(user.Profile.SocialLinks) > 0 {
			foundUser.Profile.SocialLinks = user.Profile.SocialLinks
		}
	}

	foundUser.UpdatedAt = time.Now()
	
	updatedUserModel, err := EntityToModel(&foundUser)
	if err != nil {
		log.Printf("error while converting entity to model %v", err.Error())
		return errorr.ErrInternalServer
	}

	update := bson.M{"$set": updatedUserModel}

	_, err = ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("error while update user data %v", err.Error())
		return errorr.ErrInternalServer
	}

	return nil
}

func IsUserProfileEmpty(profile entities.UserProfile) bool {
	return profile.Bio == "" &&
		profile.PictureURL == "" &&
		len(profile.SocialLinks) == 0
}
