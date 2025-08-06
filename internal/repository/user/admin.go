package userrepo

import (
	errorr "anchor-blog/internal/errors"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ur *userRepository) SetRole(ctx context.Context, id string, role string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error when cast id to object id %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	filter := bson.M{"id": objID}
	var foundUser User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error when find user by object id %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	foundUser.Role = role
	_, err = ur.collection.UpdateOne(ctx, filter, foundUser)
	if err != nil {
		log.Printf("error when update user data %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	return err
}

func (ur *userRepository) ActivateUserByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error when cast id to object id %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	filter := bson.M{"id": objID}
	var foundUser User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error when find user data %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	foundUser.Activated = true
	_, err = ur.collection.UpdateOne(ctx, filter, foundUser)
	if err != nil {
		log.Printf("error when update user data %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	return err
}

func (ur *userRepository) DeactivateUserByID(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"id": objID}
	var foundUser User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error when find user %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	foundUser.Activated = false
	_, err = ur.collection.UpdateOne(ctx, filter, foundUser)
	if err != nil {
		log.Printf("error when update user data %v \n", err.Error())
		return errorr.ErrInternalServer
	}
	return err
}
