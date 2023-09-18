package routes

import (
	"theatreManagementApp/auth"
	"theatreManagementApp/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(c *gin.Engine) {
	User := c.Group("/user")
	{
		User.POST("/signup", controller.UserSignUp)
		User.POST("/signup-verification", controller.SignupVerification)
		User.POST("/login", controller.Userlogin)
		User.GET("/profile", auth.UserAuth, controller.UserProfile)
	}
}
