package db

import (
	models "admin_user_login/Models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

var UserList []models.User

func InitDatabase() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env", err)
	}
	Db, err = gorm.Open(postgres.Open(os.Getenv("dsn")), &gorm.Config{})
	if err != nil {
		log.Fatal("Error loading database", err)
		return
	}
	err = Db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Println("Error in automigrating", err)
		return
	}

	fmt.Println("Database connected successfully")
}
