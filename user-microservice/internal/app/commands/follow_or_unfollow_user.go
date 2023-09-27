package commands

import (
	"errors"
	"music-sharing/user-microservice/internal/app/models"
	config "music-sharing/user-microservice/pkg"

	"github.com/google/uuid"
)

type FollowOrUnfollowUserCommand struct {
	Follow      bool
	UserId      uuid.UUID
	CurrentUser *models.User
}

func (cmd *FollowOrUnfollowUserCommand) Handle() error {

	db := config.Container.Database

	user := &models.User{}

	res := db.Find(user, "ID = ?", cmd.UserId)

	if res.Error != nil {
		return res.Error
	}

	if user.ID == uuid.Nil {
		return errors.New("user doesnt exist")
	}

	if cmd.Follow {
		user.Followers++
		cmd.CurrentUser.Followings++
	} else {
		user.Followers--
		cmd.CurrentUser.Followings--
	}

	db.Save(user)
	db.Save(cmd.CurrentUser)

	return nil
}
