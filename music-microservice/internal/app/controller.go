package app

import (
	"context"
	"errors"
	"mime/multipart"
	"music-sharing/music-microservice/internal/app/models"
	"music-sharing/music-microservice/internal/database"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	MusicsController struct{}

	UploadMusicReq struct {
		Title     string                `json:"title" validate:"required"`
		ShortDesc string                `json:"shortDesc" validate:"required"`
		File      *multipart.FileHeader `json:"file"`
	}

	UpdateMusicMetadataReq struct {
		Title     string `json:"title"`
		ShortDesc string `json:"shortDesc"`
	}

	RetrieveMusicsByIdsRequest struct {
		MusicsIds []string `json:"musicsIds"`
	}
)

var musicsCollection *mongo.Collection = database.OpenCollection("musics")

func (ctrl *MusicsController) GetMusics(c *gin.Context) {
	ctx := context.TODO()

	result, err := musicsCollection.Find(ctx, bson.M{})

	if err != nil {
		c.Error(err)
	}

	musics := []bson.M{}

	if err := result.All(ctx, &musics); err != nil {
		c.Error(err)
	}

	c.JSON(200, musics)
}

func (ctrl *MusicsController) GetMusicById(c *gin.Context) {
	musicId := c.Param("music_id")
	var music models.Music
	id, err := primitive.ObjectIDFromHex(musicId)

	if err != nil {
		c.Error(err)
	}

	err = musicsCollection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&music)

	if err != nil {
		c.Error(err)
	}

	c.JSON(200, music)
}

func (ctrl *MusicsController) LikeMusic(c *gin.Context) {
	musicId := c.Param("music_id")
	id, err := primitive.ObjectIDFromHex(musicId)

	if err != nil {
		c.Error(err)
	}

	var music models.Music
	filter := bson.M{"_id": id}

	err = musicsCollection.FindOne(context.TODO(), filter).Decode(&music)

	if err != nil {
		c.Error(err)
	}

	music.Likes++

	_, err = musicsCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"likes": music.Likes}})

	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func (ctrl *MusicsController) UnlikeMusic(c *gin.Context) {
	musicId := c.Param("music_id")
	id, err := primitive.ObjectIDFromHex(musicId)

	if err != nil {
		c.Error(err)
	}

	var music models.Music
	filter := bson.M{"_id": id}

	err = musicsCollection.FindOne(context.TODO(), filter).Decode(&music)

	if err != nil {
		c.Error(err)
	}

	music.Likes--

	_, err = musicsCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"likes": music.Likes}})

	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func (ctrl *MusicsController) RetrieveMusicsByIds(c *gin.Context) {
	req := RetrieveMusicsByIdsRequest{}
	var musics []models.Music

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
	}

	for _, id := range req.MusicsIds {
		var music models.Music
		parsedId, err := primitive.ObjectIDFromHex(id)

		if err != nil {
			c.Error(err)
		}

		err = musicsCollection.FindOne(context.TODO(), bson.M{"_id": parsedId}).Decode(&music)

		if err != nil {
			c.Error(err)
		}

		musics = append(musics, music)
	}

	c.JSON(200, musics)

}

func (ctrl *MusicsController) UploadMusic(c *gin.Context) {
	title := c.PostForm("title")
	shortDesc := c.PostForm("shortDesc")
	file, err := c.FormFile("file")
	userClaims := c.MustGet("user").(jwt.MapClaims)

	if err != nil {
		c.Error(err)
	}

	req := &UploadMusicReq{
		Title:     title,
		ShortDesc: shortDesc,
		File:      file,
	}

	if err := validator.New().Struct(req); err != nil {
		c.Error(err)
	}

	cid, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		c.Error(err)
	}

	res, err := cid.Upload.Upload(context.TODO(), file, uploader.UploadParams{
		PublicID: file.Filename,
	})

	if err != nil {
		c.Error(err)
	}

	music := models.Music{
		ID:        primitive.NewObjectID(),
		Title:     req.Title,
		ShortDesc: req.ShortDesc,
		PosterUrl: "",
		FileUrl:   res.SecureURL,
		Likes:     0,
		ArtistID:  userClaims["userId"].(string),
	}

	_, err = musicsCollection.InsertOne(context.TODO(), music)

	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}

func (ctrl *MusicsController) UpdateMusicMetadata(c *gin.Context) {
	musicId := c.Param("musicId")

	if len(musicId) == 0 {
		c.Error(errors.New("no musicId param found in the incoming request params"))
	}

	id, err := primitive.ObjectIDFromHex(musicId)

	if err != nil {
		c.Error(err)
	}

	req := UpdateMusicMetadataReq{}
	filter := bson.M{"_id": id}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
	}

	_, err = musicsCollection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"title": req.Title, "shortDesc": req.ShortDesc}})

	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{
		"success": true,
	})
}

func (ctrl *MusicsController) ChangeMusicPoster(c *gin.Context) {
	poster, err := c.FormFile("poster")

	if err != nil {
		c.Error(err)
	}

	musicId := c.Param("musicId")

	if len(musicId) == 0 {
		c.Error(errors.New("no musicId param found in the incoming request params"))
	}

	id, err := primitive.ObjectIDFromHex(musicId)

	if err != nil {
		c.Error(err)
	}

	cid, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))

	if err != nil {
		c.Error(err)
	}

	res, err := cid.Upload.Upload(context.TODO(), poster, uploader.UploadParams{
		PublicID: poster.Filename,
	})

	if err != nil {
		c.Error(err)
	}

	_, err = musicsCollection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"posterUrl": res.SecureURL}})

	if err != nil {
		c.Error(err)
	}

	c.JSON(200, gin.H{
		"success": true,
	})

}
