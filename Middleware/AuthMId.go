package middleware

import (
	
	"net/http"
	"strings"
	"os"
	"Backend/Config"
	"Backend/Models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func VerifToken()gin.HandlerFunc{
	return func (c *gin.Context)  {
		tokenString := c.GetHeader("Authorization")
		if tokenString == ""{
			c.JSON(http.StatusUnauthorized,gin.H{"message":"No Token provide"})
			c.Abort()
			return
		}
		if strings.HasPrefix(tokenString,"Bearer"){
			tokenString = strings.TrimPrefix(tokenString,"Bearer")
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret,nil
		})

		if err != nil || !token.Valid{
			c.JSON(http.StatusForbidden,gin.H{"message":"Invalid or expired token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ! ok || claims["email"] == nil{
			c.JSON(http.StatusForbidden,gin.H{"message":"invalid token claims"})
			c.Abort()
			return
		}
		var user models.User
		if err := config.DB.Where("email = ?",claims["email"]).First(&user).Error; err != nil{
			c.JSON(http.StatusUnauthorized,gin.H{"message":"user not found"})
			c.Abort()
			return
		}
		c.Set("user",user)
		c.Next()
	}
}