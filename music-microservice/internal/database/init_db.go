package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initDB() *mongo.Client {
	cwd, _ := os.Getwd()

	godotenv.Load(filepath.Join(cwd, "..", ".env"))

	mongoUri := os.Getenv("MONGO_URI")

	log.Printf("Connecting to %s", mongoUri)

	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))

	if err != nil {
		log.Fatal(err)
	}

	log.Print("Connected to the mongo database 🚀")

	return client
}

var Client *mongo.Client = initDB()

func OpenCollection(collectionName string) *mongo.Collection {
	return Client.Database("music-sharing").Collection(collectionName)
}
