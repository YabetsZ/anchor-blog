package userrepo

import (
	"anchor-blog/internal/domain/entities"

	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) entities.IUserRepository {
	return &userRepository{collection: collection}
}
