package mapper

import (
	entities "anchor-blog/internal/domain/entities"
	"anchor-blog/internal/dto"
	"anchor-blog/internal/infra/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ModelToEntity(model models.User) entities.User {
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

func EntityToModel(ue entities.User) (models.User, error) {
	id, err := primitive.ObjectIDFromHex(ue.ID)
	if err != nil {
		return models.User{}, err
	}
	updatedBy, err := primitive.ObjectIDFromHex(ue.UpdatedBy)
	if err != nil {
		return models.User{}, nil
	}

	userPosts := make([]primitive.ObjectID, len(ue.UserPosts))
	for index, post := range ue.UserPosts {
		userPosts[index], err = primitive.ObjectIDFromHex(post)
		if err != nil {
			return models.User{}, nil
		}
	}
	socialLinks := make([]models.SocialLink, len(ue.Profile.SocialLinks))
	for index, socialLink := range ue.Profile.SocialLinks {
		socialLinks[index] = models.SocialLink{Platform: socialLink.Platform, URL: socialLink.URL}
	}
	return models.User{
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
		Profile: models.UserProfile{
			Bio:        ue.Profile.Bio,
			PictureURL: ue.Profile.PictureURL,
			SocialLink: socialLinks,
		},
	}, nil
}

func EntityToDTO(ue entities.User) dto.UserDTO {
	socialLinks := make([]dto.SocialLinkDTO, len(ue.Profile.SocialLinks))

	for index, socialLink := range ue.Profile.SocialLinks {
		socialLinks[index] = dto.SocialLinkDTO{Platform: socialLink.Platform, URL: socialLink.URL}
	}

	return dto.UserDTO{
		ID:        ue.ID,
		Username:  ue.Username,
		FirstName: ue.FirstName,
		LastName:  ue.LastName,
		Email:     ue.Email,
		Role:      ue.Role,
		Activated: ue.Activated,
		LastSeen:  ue.LastSeen,
		CreatedAt: ue.CreatedAt,
		UpdatedBy: ue.UpdatedBy,
		UpdatedAt: ue.UpdatedAt,
		UserPosts: ue.UserPosts,
		Profile: dto.UserProfileDTO{
			Bio:         ue.Profile.Bio,
			PictureURL:  ue.Profile.PictureURL,
			SocialLinks: socialLinks,
		},
	}
}
