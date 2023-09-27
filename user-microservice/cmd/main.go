package main

import (
	"music-sharing/user-microservice/internal/app"
	"music-sharing/user-microservice/internal/app/middlewares"
	config "music-sharing/user-microservice/pkg"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	cwd, _ := os.Getwd()
	godotenv.Load(filepath.Join(cwd, "..", ".env"))

	err := config.Initialize()

	if err != nil {
		panic(err)
	}

	router := gin.Default()
	userController := &app.UserController{}

	router.Static("/profiles", "./internal/static/profiles/")

	router.Use(middlewares.ErrorHandlerMiddleware)

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"isActive": true,
		})
	})

	router.POST("/login", userController.Login)
	router.POST("/register", userController.Register)
	router.GET("/myProfile", middlewares.AuthMiddleware, userController.MyProfile)
	router.GET("/viewProfile/:userId", userController.ViewProfile)
	router.GET("/followUser/:userId", middlewares.AuthMiddleware, userController.FollowUser)
	router.GET("/unfollowUser/:userId", middlewares.AuthMiddleware, userController.UnfollowUser)
	router.GET("/updateMyAccount", middlewares.AuthMiddleware, userController.UpdateMyAccount)
	router.POST("/uploadProfile", middlewares.AuthMiddleware, userController.UploadProfile)

	router.Run()
}
