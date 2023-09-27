package queries

import (
	"errors"
	"music-sharing/user-microservice/internal/app/models"
	config "music-sharing/user-microservice/pkg"

	"github.com/google/uuid"
)

type (
	GetUserProfileByIdQueryResponse struct {
		User *models.User `json:"user"`
	}

	GetUserProfileByIdQuery struct {
		ID uuid.UUID
	}
)

func (c *GetUserProfileByIdQuery) Handle() (interface{}, error) {

	db := config.Container.Database

	user := &models.User{}

	res := db.Find(user, "ID = ?", c.ID)

	if user.ID == uuid.Nil {
		return nil, errors.New("user doesnt exist")
	}

	if res.Error != nil {
		return nil, res.Error
	}

	resp := &GetUserProfileByIdQueryResponse{
		User: user,
	}

	return resp, nil

}