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
		City.GET("", controller.GetCityList)
		City.POST("/add", auth.AdminAuth, controller.AddCity)
		City.PUT("/edit/:id", auth.AdminAuth, controller.EditCity)
		City.DELETE("/delete/:id", auth.AdminAuth, controller.DeleteCity)

		//Cinemas Management
		Cinemas := Admin.Group("/cinemas")
		Cinemas.POST("/add", auth.AdminAuth, controller.AddCinemas)
		Cinemas.PUT("/edit/:id", auth.AdminAuth, controller.EditCinemas)
		Cinemas.DELETE("/delete/:id", auth.AdminAuth, controller.DeleteCinemas)

		//Manager role management
		Manager := Admin.Group("/manager")
		Manager.POST("/add", auth.AdminAuth, controller.AddManager)
		Manager.PUT("/edit/:id", auth.AdminAuth, controller.EditManager)
		Manager.DELETE("/delete/:id", auth.AdminAuth, controller.DeleteManager)
	}
}
