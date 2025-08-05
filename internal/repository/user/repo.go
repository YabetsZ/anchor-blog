package userrepo

import "go.mongodb.org/mongo-driver/mongo"

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *userRepository {
	return &userRepository{collection: collection}
}
