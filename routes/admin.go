package routes

import (
	"theatreManagementApp/controller"

	"github.com/gin-gonic/gin"
)

func AdminRoutes(c *gin.Engine) {
	Admin := c.Group("/admin")
	{
		Admin.POST("/login", controller.AdminLogin)
	}
}
