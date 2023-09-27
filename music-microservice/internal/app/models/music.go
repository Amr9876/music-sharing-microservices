package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Music struct {
	ID        primitive.ObjectID `bson:"_id"`
	ArtistID  string             `bson:"artistId" json:"artistId"`
	Likes     uint               `json:"likes"`
	FileUrl   string             `json:"fileUrl"`
	PosterUrl string             `json:"posterUrl"`
	Title     string             `json:"title"`
	ShortDesc string             `json:"shortDesc"`
}
