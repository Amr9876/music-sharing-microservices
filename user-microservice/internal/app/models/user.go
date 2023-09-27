package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID             uuid.UUID `json:"id" gorm:"primaryKey"`
	FullName       string    `json:"fullName"`
	Gender         string    `json:"gender"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashedPassword"`
	Followers      uint      `json:"followers"`
	Followings     uint      `json:"followings"`
	ProfileURL     string    `json:"profileUrl"`
	IsPrivate      bool      `json:"isPrivate"`
}
