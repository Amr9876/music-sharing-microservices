package lib

import (
	"context"
	"log"
	"math/rand"
	"music-sharing/music-microservice/internal/app/models"
	"music-sharing/music-microservice/internal/database"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Seeds 10 random musics to the mongo database if it's empty
func SeedMusics() error {

	musicsCollection := database.OpenCollection("musics")

	count, err := musicsCollection.CountDocuments(context.TODO(), bson.M{})

	if err != nil {
		return err
	}

	// exists the func if the collection arleady have documents
	if count != 0 {
		return nil
	}

	gofakeit.Seed(0)

	for i := 0; i < 10; i++ {
		music := models.Music{
			ID:        primitive.NewObjectID(),
			Likes:     uint(rand.Int()),
			Title:     gofakeit.HipsterWord(),
			ShortDesc: gofakeit.HipsterSentence(15),
			PosterUrl: gofakeit.ImageURL(250, 250),
			FileUrl:   "https://res.cloudinary.com/dpoxxjpmu/video/upload/v1688663505/l55gmkryd9u82kql09no.mp3",
			ArtistID:  "b3828065-44e9-4923-8e82-6ca03998a6c4",
		}

		_, err := musicsCollection.InsertOne(context.TODO(), music)

		if err != nil {
			return err
		}
	}

	log.Print("No documents found 😕, so seed 10 musics successfully 🔥")

	return nil

}
