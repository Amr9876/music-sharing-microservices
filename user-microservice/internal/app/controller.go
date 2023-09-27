package app

import (
	"errors"
	"io"
	"music-sharing/user-microservice/internal/app/commands"
	"music-sharing/user-microservice/internal/app/models"
	"music-sharing/user-microservice/internal/app/queries"
	"music-sharing/user-microservice/internal/lib"
	config "music-sharing/user-microservice/pkg"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type (
	LoginBody struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	RegisterBody struct {
		FullName  string `json:"fullName" validate:"required"`
		Gender    string `json:"gender" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
		IsPrivate bool   `json:"isPrivate"`
	}

	UpdateAccountBody struct {
		ID        uuid.UUID `json:"id" gorm:"primaryKey"`
		FullName  string    `json:"fullName"`
		Email     string    `json:"email"`
		IsPrivate bool      `json:"isPrivate"`
	}

	UserController struct{}
)

func (ctrl *UserController) Login(c *gin.Context) {

	db := config.Container.Database
	loginBody := &LoginBody{}
	user := &models.User{}

	err := lib.BindAndValidate(loginBody, c)

	if err != nil {
		c.Error(err)
		return
	}

	resp := db.Find(user, "Email = ?", loginBody.Email)

	if resp.Error != nil {
		c.Error(resp.Error)
		return
	}

	if user.ID == uuid.Nil {
		c.Error(errors.New("user doesnt exist"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginBody.Password)); err != nil {
		c.Error(err)
		return
	}

	tokenString, err := lib.CreateJWT(map[string]interface{}{
		"userId":    user.ID.String(),
		"isPrivate": user.IsPrivate,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"token": tokenString,
	})

}

func (ctrl *UserController) Register(c *gin.Context) {

	body := &RegisterBody{}
	commandBus := config.Container.CommmandBus
	err := lib.BindAndValidate(body, c)

	if err != nil {
		c.Error(err)
		return
	}

	respErr := commandBus.Send(&commands.CreateUserCommand{
		FullName:  body.FullName,
		Gender:    body.Gender,
		Email:     body.Email,
		Password:  body.Password,
		IsPrivate: body.IsPrivate,
	})

	if respErr != nil {
		c.Error(respErr)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}

func (ctrl *UserController) MyProfile(c *gin.Context) {

	user, exists := c.Get("user")

	if !exists {
		c.Error(errors.New("user doesnt exist"))
		return
	}

	c.JSON(200, user.(*models.User))

}

func (ctrl *UserController) ViewProfile(c *gin.Context) {

	userId := c.Param("userId")
	queryBus := config.Container.QueryBus

	resp, err := queryBus.Send(&queries.GetUserProfileByIdQuery{
		ID: uuid.MustParse(userId),
	})

	if err != nil {
		c.Error(err)
		return
	}

	user := resp.(*queries.GetUserProfileByIdQueryResponse).User

	if user.IsPrivate {
		c.Error(errors.New("this account is private"))
		return
	}

	c.JSON(200, user)

}

func (ctrl *UserController) FollowUser(c *gin.Context) {
	userId := c.Param("userId")
	user := c.MustGet("user").(*models.User)
	commandBus := config.Container.CommmandBus

	if len(userId) == 0 {
		c.Error(errors.New("user id is required"))
		return
	}

	err := commandBus.Send(&commands.FollowOrUnfollowUserCommand{
		Follow:      true,
		UserId:      uuid.MustParse(userId),
		CurrentUser: user,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}

func (ctrl *UserController) UnfollowUser(c *gin.Context) {
	userId := c.Param("userId")
	user := c.MustGet("user").(*models.User)
	commandBus := config.Container.CommmandBus

	if len(userId) == 0 {
		c.Error(errors.New("user id is required"))
		return
	}

	err := commandBus.Send(&commands.FollowOrUnfollowUserCommand{
		Follow:      false,
		UserId:      uuid.MustParse(userId),
		CurrentUser: user,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func (ctrl *UserController) UpdateMyAccount(c *gin.Context) {

	db := config.Container.Database
	user := c.MustGet("user").(*models.User)
	body := &UpdateAccountBody{}

	err := lib.BindAndValidate(body, c)

	if err != nil {
		c.Error(err)
		return
	}

	if user.ID != body.ID {
		c.Error(errors.New("user id doesnt match"))
		return
	}

	resp := db.Save(user)

	if resp.Error != nil {
		c.Error(resp.Error)
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"user":    user,
	})

}

func (ctrl *UserController) UploadProfile(c *gin.Context) {

	db := config.Container.Database
	user := c.MustGet("user").(*models.User)

	file, header, err := c.Request.FormFile("profile")

	if err != nil {
		c.Error(err)
		return
	}

	fileUrl := "./internal/static/profiles/" + header.Filename

	out, err := os.Create(fileUrl)

	if err != nil {
		c.Error(err)
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)

	if err != nil {
		c.Error(err)
		return
	}

	user.ProfileURL = "/profiles/" + header.Filename

	db.Save(user)

	c.JSON(200, gin.H{
		"success": true,
	})
}
