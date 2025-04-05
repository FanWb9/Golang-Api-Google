package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB * gorm.DB

func Connect(){
	_ = godotenv.Load()

	dsn := os.Getenv("DB")
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal("Database Tidak connect:", err)
	}
	DB = database
	log.Println("Datatbase Connect")
}