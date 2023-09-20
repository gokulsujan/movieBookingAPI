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
		City.GET("", auth.AdminAuth, controller.GetCityList)
		City.POST("/add", auth.AdminAuth, controller.AddCity)
		City.PUT("/edit/:id", auth.AdminAuth, controller.EditCity)
		City.DELETE("/delete/:id", auth.AdminAuth, controller.DeleteCity)

		//Cinemas Management
		Cinemas := Admin.Group("/cinemas")
		Cinemas.GET("/:id", auth.AdminAuth, controller.GetCinemasListFromCities)
		Cinemas.POST("/add", auth.AdminAuth, controller.AddCinemas)
		Cinemas.PUT("/edit/:id", auth.AdminAuth, controller.EditCinemas)
		Cinemas.DELETE("/delete/:id", auth.AdminAuth, controller.DeleteCinemas)

		//Manager role management
		Manager := Admin.Group("/manager")
		Manager.GET("", auth.AdminAuth, controller.GetManagerList)
		Manager.POST("/add", auth.AdminAuth, controller.AddManager)
		Manager.PUT("/edit/:id", auth.AdminAuth, controller.EditManager)
		Manager.DELETE("/delete/:id", auth.AdminAuth, controller.DeleteManager)

		//Movie Management
		Movies := Admin.Group("/movies")
		Movies.POST("/add", auth.AdminAuth, controller.AddMovies)
		Movies.PUT("/edit/:id", auth.AdminAuth, controller.EditMovies)
		Movies.DELETE("/delete/:id", auth.AdminAuth, controller.DeleteMovies)

		//User Management
		Users := Admin.Group("/users")
		Users.GET("", auth.AdminAuth, controller.GetUsersList)
		Users.PUT("/status/:id", auth.AdminAuth, controller.UserStatusChange)

		//screenformat
		ScreenFromat := Admin.Group("/screen-format")
		ScreenFromat.GET("", auth.AdminAuth, controller.ViewScreenFormat)
		ScreenFromat.POST("/add", auth.AdminAuth, controller.AddScreenFormat)
	}
}
