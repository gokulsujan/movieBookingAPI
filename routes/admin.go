package routes

import (
	"theatreManagementApp/auth"
	"theatreManagementApp/controller"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(c *gin.Engine) {
	Admin := c.Group("/admin")
	{
		Admin.POST("/login", controller.AdminLogin)

		//City Management
		City := Admin.Group("/city")
		City.GET("")
		City.POST("/add", auth.AdminAuth, controller.AddCity)
	}
}
