package middlewares

import (
	"errors"
	"music-sharing/music-microservice/internal/lib"
	"strings"

	"github.com/gin-gonic/gin"
)

func IsOwnerMiddleware(c *gin.Context) {
	header := c.Request.Header.Get("Authorization")

	if len(header) == 0 {
		c.Error(errors.New("no authorization header"))
	}

	ownerId := c.Param("ownerId")

	if len(ownerId) == 0 {
		c.Error(errors.New("no owner_id param found in the incoming request params"))
	}

	tokenString := strings.Split(header, " ")[1]

	userClaims, err := lib.ParseJWT(tokenString)

	if err != nil {
		c.Error(err)
	}

	if userClaims["userId"] != ownerId {
		c.Error(errors.New("you are not authorized"))
	} else {
		c.Next()
	}
}
