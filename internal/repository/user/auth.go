package userrepo

import (
	"context"
	"log"

	errorr "anchor-blog/internal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	type IUserAuthRepository interface {
		CheckEmail(ctx context.Context, email string) (bool, error)
		CheckUsername(ctx context.Context, username string) (bool, error)
		ChangePassword(ctx context.Context, id string, newHashedPassword string) error
		ChangeEmail(ctx context.Context, email string, newEmail string) error
	}
*/
func (ur *userRepository) CheckEmail(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := ur.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("error when find email %v \n", err.Error())
		return true, errorr.ErrInternalServer
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}

func (ur *userRepository) CheckUsername(ctx context.Context, username string) (bool, error) {
	filter := bson.M{"username": username}
	count, err := ur.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("error when find username %v \n", err.Error())
		return true, errorr.ErrInternalServer
	}
	if count == 0 {
		return false, err
	}
	return true, err
}

func (ur *userRepository) ChangePassword(ctx context.Context, id string, newPassword string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error when cast id to object id %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	filter := bson.M{"_id": objID}

	var foundUser User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error when find by object id %v \n", err.Error())
		return errorr.ErrInternalServer
	}

	foundUser.PasswordHash = newPassword

	_, err = ur.collection.UpdateOne(ctx, filter, foundUser)
	if err != nil {
		log.Printf("error when update user data %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	return nil
}

func (ur *userRepository) ChangeEmail(ctx context.Context, email string, newEmail string) error {

	filter := bson.M{"email": email}

	var foundUser User
	err := ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error when find user %v \n", err.Error())
		return errorr.ErrInternalServer
	}

	foundUser.Email = newEmail
	_, err = ur.collection.UpdateOne(ctx, filter, foundUser)
	if err != nil {
		log.Printf("error when update user data %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	return nil
}
