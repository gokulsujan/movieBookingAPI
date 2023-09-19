package routes

import (
	"theatreManagementApp/auth"
	"theatreManagementApp/controller"

	"github.com/gin-gonic/gin"
)

func ManagerRoutes(c *gin.Engine) {
	Manager := c.Group("/manager")
	{
		Manager.POST("/login", controller.ManagerLogin)

		//Screen Management
		screen := Manager.Group("/screen")
		screen.POST("/add", auth.ManagerAuth, controller.AddScreen)
	}
}
