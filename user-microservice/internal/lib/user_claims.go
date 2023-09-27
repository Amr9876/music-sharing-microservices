package lib

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	UserId    string `json:"userId"`
	IsPrivate bool   `json:"isPrivate"`
	jwt.Claims
}
