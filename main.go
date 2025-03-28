package main

import (
	"Backend/Config"
	"Backend/Routes"
	"log"
	

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)
func main(){
	
	err := godotenv.Load()
	if err != nil{
		log.Fatal("eror load env")
	}
	
	r := gin.Default()
	config.Connect()
	routes.AuthRoutes(r)
	r.Run(":8080") 


	
}