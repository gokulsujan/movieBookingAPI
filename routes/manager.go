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
		Manager.POST("/forget-password", controller.ManagerForgetPass)
		Manager.POST("/change-password", controller.ChangePassword)

		//Screen Management
		screen := Manager.Group("/screen")
		screen.GET("", auth.ManagerAuth, controller.GetScreenList)
		screen.POST("/add", auth.ManagerAuth, controller.AddScreen)
		screen.PUT("/edit/:id", auth.ManagerAuth, controller.EditScreen)
		screen.DELETE("/delete/:id", auth.ManagerAuth, controller.DeleteScreen)

		//Shows Management
		shows := Manager.Group("/shows")
		shows.GET("", auth.ManagerAuth, controller.GetRunnigMovies)
		shows.POST("/add", auth.ManagerAuth, controller.AddShows)
		shows.PUT("/change_status", auth.ManagerAuth, controller.ShowStatusChange)

	}
}
