package infrastructure

import (
	"log"
	"music-sharing/user-microservice/internal/app/models"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {

	db, _ = gorm.Open(mysql.Open(os.Getenv("MYSQL_CONN")), &gorm.Config{})

	db.AutoMigrate(&models.User{})

	log.Printf("🚀 Connected to %s", os.Getenv("MYSQL_CONN"))

	return db
}
