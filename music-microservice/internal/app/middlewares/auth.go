package middlewares

import (
	"errors"
	"music-sharing/music-microservice/internal/lib"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {

	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		c.Error(errors.New("no authorization header"))
	}

	tokenString := strings.Split(header, " ")[1]

	userClaims, err := lib.ParseJWT(tokenString)

	if err != nil {
		c.Error(err)
	}

	c.Set("user", userClaims)
	c.Set("user_token", tokenString)

	c.Next()
}
