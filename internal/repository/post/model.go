package postrepo

import (
	"anchor-blog/internal/domain/entities"
	"anchor-blog/internal/errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title     string               `bson:"title" json:"title"`
	Content   string               `bson:"content" json:"content"`
	AuthorID  primitive.ObjectID   `bson:"author_id" json:"author_id"`
	Tags      []string             `bson:"tags" json:"tags"`
	ViewCount int                  `bson:"view_count" json:"view_count"`
	Likes     []primitive.ObjectID `bson:"likes" json:"likes"`
	Dislikes  []primitive.ObjectID `bson:"dislikes" json:"dislikes"`
	CreatedAt time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time            `bson:"updated_at" json:"updated_at"`
}

// ::::::: Mapping functions :::::::::::
func ToDomainPost(p *Post) *entities.Post {
	return &entities.Post{
		ID:        p.ID.Hex(),
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  p.AuthorID.Hex(),
		Tags:      p.Tags,
		ViewCount: p.ViewCount,
		Likes:     objectIDsToHex(p.Likes),
		Dislikes:  objectIDsToHex(p.Dislikes),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func objectIDsToHex(ids []primitive.ObjectID) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = id.Hex()
	}
	return result
}

func FromDomainPost(p *entities.Post) (*Post, error) {
	id, err := primitive.ObjectIDFromHex(p.ID)
	if err != nil {
		return nil, err
	}
	authorID, err := primitive.ObjectIDFromHex(p.AuthorID)
	if err != nil {
		return nil, err
	}
	likes, err := hexToObjectIDs(p.Likes)
	if err != nil {
		return nil, err
	}
	dislikes, err := hexToObjectIDs(p.Dislikes)
	if err != nil {
		return nil, err
	}

	return &Post{
		ID:        id,
		Title:     p.Title,
		Content:   p.Content,
		AuthorID:  authorID,
		Tags:      p.Tags,
		ViewCount: p.ViewCount,
		Likes:     likes,
		Dislikes:  dislikes,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

func hexToObjectIDs(hexIDs []string) ([]primitive.ObjectID, error) {
	objs := make([]primitive.ObjectID, len(hexIDs))
	for i, hex := range hexIDs {
		id, err := primitive.ObjectIDFromHex(hex)
		if err != nil {
			log.Println("invalid id format ", id)
			return nil, errors.ErrInvalidPostID
		}
		objs[i] = id
	}
	return objs, nil
}
