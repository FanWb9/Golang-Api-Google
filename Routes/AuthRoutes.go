package routes

import(
	"Backend/Controller"
	"github.com/gin-gonic/gin"
)

func AuthRoutes( r * gin.Engine){
	auth := r.Group("/api/auth")
	{
		auth.POST("/register",controller.Register)
		auth.POST("/login",controller.Login)
		auth.POST("/login-google",controller.LoginWithGoogle)
	}
}
