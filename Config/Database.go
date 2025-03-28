package config

import(
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"log"
)

var DB * gorm.DB

func Connect(){
	dsn := "root:@tcp(127.0.0.1:3306)/catalog_db?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil{
		log.Fatal("Database Tidak connect:", err)
	}
	DB = database
	log.Println("Datatbase Connect")
}