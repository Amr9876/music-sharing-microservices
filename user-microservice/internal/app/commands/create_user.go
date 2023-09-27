package commands

import (
	"music-sharing/user-microservice/internal/app/models"
	config "music-sharing/user-microservice/pkg"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	CreateUserCommand struct {
		FullName  string
		Gender    string
		Email     string
		Password  string
		IsPrivate bool
	}
)

func (cmd *CreateUserCommand) Handle() error {

	db := config.Container.Database

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := &models.User{
		ID:             uuid.New(),
		FullName:       cmd.FullName,
		Gender:         cmd.Gender,
		Email:          cmd.Email,
		HashedPassword: string(hashedPassword),
		Followers:      0,
		Followings:     0,
		ProfileURL:     "",
		IsPrivate:      cmd.IsPrivate,
	}

	res := db.Create(user)

	if res.Error != nil {
		return res.Error
	}

	return nil

}
