package lib

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWT(tokenString string) (jwt.MapClaims, error) {

	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil

}
