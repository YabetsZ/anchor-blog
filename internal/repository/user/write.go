package userrepo

import (
	"anchor-blog/internal/domain/entities"
	"context"
	"log"
	"time"

	AppError "anchor-blog/internal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ur *userRepository) CreateUser(ctx context.Context, user *entities.User) (string, error) {
	user.ID = primitive.NewObjectID().Hex() // Needs to be modified
	user.UpdatedBy = user.ID
	user.UserPosts = make([]string, 0)

	userDoc, err := EntityToModel(user)
	if err != nil {
		log.Printf("error while transfer user entity to user model %v", err.Error())
		return "", err
	}

	_, err = ur.collection.InsertOne(ctx, userDoc)
	if err != nil {
		log.Printf("error while create new user %v", err.Error())
		return "", AppError.ErrInternalServer
	}
	return userDoc.ID.Hex(), nil
}
func (ur *userRepository) DeleteUserByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return AppError.ErrInternalServer
	}
	filter := bson.M{"id": objID}
	_, err = ur.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("error while delete user %v", err.Error())
		return AppError.ErrInternalServer
	}
	return nil
}

func (ur *userRepository) SetLastSeen(ctx context.Context, id string, timestamp time.Time) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error while cast id to object id %v", err.Error())
		return AppError.ErrInternalServer
	}
	filter := bson.M{"id": objID}

	var foundUser User

	err = ur.collection.FindOne(ctx, filter).Decode(foundUser)
	if err != nil {
		log.Printf("error while find user %v", err.Error())
		return AppError.ErrInternalServer
	}
	foundUser.LastSeen = timestamp
	_, err = ur.collection.UpdateOne(ctx, filter, foundUser)
	if err != nil {
		return AppError.ErrInternalServer
	}
	return err
}

func (ur *userRepository) EditUserByID(ctx context.Context, id string, user *entities.User) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error while cast id to object id %v", err.Error())
		return AppError.ErrInternalServer
	}

	filter := bson.M{"_id": objID}

	var foundUserModel User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUserModel)
	if err != nil {
		log.Printf("error while find user %v", err.Error())
		return AppError.ErrInternalServer
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
		return AppError.ErrInternalServer
	}

	update := bson.M{"$set": updatedUserModel}

	_, err = ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("error while update user data %v", err.Error())
		return AppError.ErrInternalServer
	}

	return nil
}

func IsUserProfileEmpty(profile entities.UserProfile) bool {
	return profile.Bio == "" &&
		profile.PictureURL == "" &&
		len(profile.SocialLinks) == 0
}

func (ur *userRepository) UpdateUserRole(ctx context.Context, adminID, targetID, role string) error {
	adminObjID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		log.Printf("invalid promoterID %s: %v", adminID, err)
		return AppError.ErrInvalidUserID
	}
	targetObjID, err := primitive.ObjectIDFromHex(targetID)
	if err != nil {
		log.Printf("invalid targetID %s: %v", targetID, err)
		return AppError.ErrInvalidUserID
	}
	filter := bson.M{
		"_id": targetObjID,
	}
	update := bson.M{
		"$set": bson.M{
			"role":       role,
			"updated_by": adminObjID,
		},
	}
	_, err = ur.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return AppError.ErrInternalServer
	}
	return nil
}
