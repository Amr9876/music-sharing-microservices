package main

import (
	"log"
	"music-sharing/music-microservice/internal/app"
	"music-sharing/music-microservice/internal/app/middlewares"
	"music-sharing/music-microservice/internal/lib"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	cwd, _ := os.Getwd()
	err := godotenv.Load(filepath.Join(cwd, "..", ".env"))

	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")

	err = lib.SeedMusics()

	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	controller := app.MusicsController{}

	router.Use(middlewares.AuthMiddleware)
	router.Use(middlewares.ErrorHandlerMiddleware)

	router.GET("/getMusicById/:music_id", controller.GetMusicById)
	router.GET("/getMusics", controller.GetMusics)
	router.POST("/likeMusic/:music_id", controller.LikeMusic)
	router.POST("/unlikeMusic/:music_id", controller.UnlikeMusic)
	router.POST("/retrieveMusicsByIds", controller.RetrieveMusicsByIds)
	router.POST("/uploadMusic", controller.UploadMusic)
	router.POST("/updateMusicMetadata/:ownerId/:musicId", middlewares.IsOwnerMiddleware, controller.UpdateMusicMetadata)
	router.POST("/changeMusicPoster/:ownerId/:musicId", middlewares.IsOwnerMiddleware, controller.ChangeMusicPoster)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"isActive": true,
		})
	})

	router.Run(port)

}
