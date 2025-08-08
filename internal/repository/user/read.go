package userrepo

import (
	"anchor-blog/internal/domain/entities"
	"context"
	"log"

	errorr "anchor-blog/internal/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (ur *userRepository) GetUserByID(ctx context.Context, id string) (*entities.User, error) {
	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error while cast id to object id %v", err.Error())
		return &entities.User{}, errorr.ErrInternalServer
	}
	filter := bson.M{"_id": ObjID}
	var foundUser User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error while find user %v", err.Error())
		return &entities.User{}, errorr.ErrInternalServer
	}
	user := ModelToEntity(&foundUser)

	return &user, nil
}

func (ur *userRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	filter := bson.M{"email": email}
	var foundUser User
	err := ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error while find user %v", err.Error())
		return &entities.User{}, errorr.ErrInternalServer
	}
	user := ModelToEntity(&foundUser)
	return &user, nil
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*entities.User, error) {
	filter := bson.M{"username": username}
	var foundUser User
	err := ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error while find user %v", err.Error())
		return &entities.User{}, errorr.ErrInternalServer
	}
	user := ModelToEntity(&foundUser)
	return &user, nil
}

func (ur *userRepository) GetUsers(ctx context.Context, limit, offset int64) ([]*entities.User, error) {
	// pagination
	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(offset)

	filter := bson.M{}
	cursor, err := ur.collection.Find(ctx, filter, opts)

	if err != nil {
		log.Printf("error while find cursor %v", err.Error())
		return nil, errorr.ErrInternalServer
	}
	defer cursor.Close(ctx)

	var users []*entities.User
	for cursor.Next(ctx) {
		var userDoc User
		if err := cursor.Decode(&userDoc); err != nil {
			log.Printf("error while decode user to address %v", err.Error())
			return nil, errorr.ErrInternalServer
		}
		user := ModelToEntity(&userDoc)
		users = append(users, &user)

	}
	if err := cursor.Err(); err != nil {
		log.Printf("error while check cursor %v", err.Error())
		return nil, errorr.ErrInternalServer
	}
	return users, nil

}

func (ur *userRepository) CountUsersByRole(ctx context.Context, role string) (int64, error) {
	filter := bson.M{"role": role}
	count, err := ur.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("error while count user by role %v", err.Error())
		return 0, errorr.ErrInternalServer
	}
	return count, nil
}

func (ur *userRepository) CountActiveUsers(ctx context.Context) (int64, error) {
	filter := bson.M{"activated": true}
	count, err := ur.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("error while count active users %v", err.Error())
		return 0, errorr.ErrInternalServer
	}
	return count, nil
}

func (ur *userRepository) CountAllUsers(ctx context.Context) (int64, error) {
	filter := bson.M{}
	count, err := ur.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("error while count users %v", err.Error())
		return 0, errorr.ErrInternalServer
	}
	return count, nil
}

func (ur *userRepository) CountInactiveUsers(ctx context.Context) (int64, error) {
	filter := bson.M{"activated": false}
	count, err := ur.collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Printf("error while count inactive users %v", err.Error())
		return 0, errorr.ErrInternalServer
	}
	return count, nil
}

func (ur *userRepository) GetUserRoleByID(ctx context.Context, id string) (string, error) {
	ObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("error while cast id to object id %v", err.Error())
		return "", errorr.ErrInternalServer
	}
	filter := bson.M{"id": ObjID}
	var foundUser User
	err = ur.collection.FindOne(ctx, filter).Decode(&foundUser)
	if err != nil {
		log.Printf("error while count user  %v", err.Error())
		return "", errorr.ErrInternalServer
	}

	return foundUser.Role, nil
}
