package queries

import (
	"errors"
	"music-sharing/user-microservice/internal/app/models"
	"music-sharing/user-microservice/internal/lib"
	config "music-sharing/user-microservice/pkg"
	"strings"

	"github.com/google/uuid"
)

type (
	GetUserProfileByTokenQueryResponse struct {
		User *models.User `json:"user"`
	}

	GetUserProfileByTokenQuery struct {
		Token string
	}
)

func (c *GetUserProfileByTokenQuery) Handle() (interface{}, error) {

	db := config.Container.Database

	user := &models.User{}
	claims, err := lib.ParseJWT(c.Token)
	userId := claims["userId"].(string)

	if err != nil {
		return nil, err
	}

	res := db.Find(user, "id = ?", userId)

	if user.ID == uuid.Nil {
		return nil, errors.New(strings.Join([]string{"tried to fetch user but doesnt exist", userId}, " "))
	}

	if res.Error != nil {
		return nil, res.Error
	}

	resp := &GetUserProfileByTokenQueryResponse{
		User: user,
	}

	return resp, nil

}
