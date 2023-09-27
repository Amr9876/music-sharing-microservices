package middlewares

import (
	"errors"
	"fmt"
	"music-sharing/user-microservice/internal/app/queries"
	config "music-sharing/user-microservice/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {

	queryBus := config.Container.QueryBus
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		c.Error(errors.New("no authorization header"))
		return
	}

	tokenString := strings.Split(header, " ")[1]

	resp, err := queryBus.Send(&queries.GetUserProfileByTokenQuery{
		Token: tokenString,
	})

	if err != nil {
		c.Error(err)
		return
	}

	user := resp.(*queries.GetUserProfileByTokenQueryResponse).User
	fmt.Println(user)
	c.Set("user", user)

	c.Next()
}
