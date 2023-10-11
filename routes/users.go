package routes

import (
	"theatreManagementApp/auth"
	"theatreManagementApp/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(c *gin.Engine) {
	User := c.Group("")
	{
		User.POST("/signup", controller.UserSignUp)
		User.POST("/signup-verification", controller.SignupVerification)
		User.POST("/login", controller.Userlogin)
		User.GET("/profile", auth.UserAuth, controller.UserProfile)
		User.GET("/bookings", auth.UserAuth, controller.UserBookings)
		User.POST("/forget-password", controller.UserPassChange)
		User.POST("/change-password", controller.UserChangePassword)
		User.GET("/home", auth.UserAuth, controller.SelectCity)
		User.GET("/cinemas/:city", auth.UserAuth, controller.CinemasList)
		User.GET("/:city", auth.UserAuth, controller.MoviesList)
		User.GET("/:city/:id", auth.UserAuth, controller.CinemasListOfMovies)
		User.GET("/shows/:cinemas/:id", auth.UserAuth, controller.ShowsListByCinemas)
		User.GET("/show/seats/:id", auth.UserAuth, controller.BookingLayout)
		User.POST("/show/booking/:id", auth.UserAuth, controller.BookSeats)
		User.GET("/show/payment/:id", controller.PaymentPage)
		User.POST("/booking/paymentSuceess", controller.PaymentValidation)
	}
}
