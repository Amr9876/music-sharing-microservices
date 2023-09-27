package lib

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func BindAndValidate(body interface{}, c *gin.Context) error {

	if err := c.ShouldBindJSON(body); err != nil {
		return err
	}

	if err := validator.New().Struct(body); err != nil {
		return err
	}

	return nil
}