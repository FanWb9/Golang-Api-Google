package controller

import (
	"Backend/Config"
	"Backend/Models"
	"net/http"
	"os"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

var googleAuthConfig = &oauth2.Config{
	ClientID : os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret : os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL: "http://localhost:8080/api/auth/google/callback",
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint: google.Endpoint,
}

func Register(c *gin.Context){
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password),10)
	user.Password = string(hashedPassword)

	if err := config.DB.Create(&user).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"message":"failed to register user"})
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"user register success"})

}

func Login (c* gin.Context){
	var input models.User
	var user models.User

	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"message":"error"})
		return
	}
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil{
		c.JSON(http.StatusUnauthorized,gin.H{"message":"invalid email or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(input.Password)); err != nil{
		c.JSON(http.StatusUnauthorized,gin.H{"message":"invalid email or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"email" : user.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{"token":tokenString})
}
func LoginWithGoogle(c *gin.Context) {
	var input struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	var user models.User

	// Debug: Cek apakah request diterima
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Gagal bind JSON:", err) // Tambahkan log error
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	fmt.Println("Request diterima - Email:", input.Email, "Name:", input.Name) // Debug

	// Cek apakah email sudah ada di database
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		user = models.User{
			Email: input.Email,
			Name:  input.Name,
		}
		// Tambahkan log sebelum menyimpan ke database
		fmt.Println("User belum ada, menambahkan ke database:", user.Email)

		if err := config.DB.Create(&user).Error; err != nil {
			fmt.Println("Gagal menyimpan user ke database:", err) // Debug error database
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed Login google"})
			return
		}
	}

	// Debug: Cek apakah user ditemukan/ditambahkan dengan benar
	fmt.Println("User ditemukan atau berhasil ditambahkan:", user.Email)

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	// Debug: Cek apakah token berhasil dibuat
	fmt.Println("Token berhasil dibuat:", tokenString)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}
