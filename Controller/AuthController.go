package controller

import (
	"Backend/Config"
	"Backend/Models"
	"net/http"
	"os"
	"time"

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
	var CheckedVald models.User //this is to check if the user has already registered
	if err := c.ShouldBindJSON(&user); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}
	//This Check is to register if user has alredy resgister via the email  and phoneNumber same
	if err := config.DB.Where("email = ? or phone_number = ?",user.Email, user.PhoneNumber).First(&CheckedVald).Error; err == nil{
		if CheckedVald.Email == user.Email{
			c.JSON(http.StatusBadRequest,gin.H{"message":"email alredy exist"})
		}else{
			c.JSON(http.StatusBadRequest,gin.H{"message":"phone number alredy exist"})
		}
		return
	}
	//this generate password to hash
	//and store it in the database
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password),10)
	user.Password = string(hashedPassword)

	//this convert date of birth to time.Time with the format(yyyy-mm-dd)
	//and store it in the database
	//this is the format of date of birth example 2006-01-02
	Layout := "2006-01-02" 
	parsedDate , err := time.Parse(Layout, user.DateBirtStr)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"message":"invalid date format"})
		return
	}
	user.DateOfBirth = &parsedDate

	//Create account for NewMember
	if err := config.DB.Create(&user).Error; err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{"message":"failed to register user"})
		return
	}
	
	c.JSON(http.StatusOK,gin.H{"message":"user register success"})

}
//This function Login 
func Login (c* gin.Context){
	//create variabel input or user with parameters models user
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request"})
		return
	}

	// Cek apakah email sudah ada di database
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		user = models.User{
			Email: input.Email,
			Name:  input.Name,
		}
		//membuat akun jika belum terdaftar di database
		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed Login google"})
			return
		}
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user":  user,
	})
}
