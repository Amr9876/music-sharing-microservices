package models

import (
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	FullName       string    `json:"fullName"`
	Gender         string    `json:"gender"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashedPassword"`
	Followers      uint      `json:"followers"`
	Followings     uint      `json:"followings"`
	ProfileURL     string    `json:"profileUrl"`
	IsPrivate      bool      `json:"isPrivate"`
}
