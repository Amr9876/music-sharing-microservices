package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware(c *gin.Context) {

	c.Next()

	errors := c.Errors

	if len(errors) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors":  errors.Errors(),
			"success": false,
		})
	}

}
