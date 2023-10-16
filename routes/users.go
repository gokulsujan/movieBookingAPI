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
		User.GET("/cinemas", auth.UserAuth, controller.CinemasList)
		User.GET("", auth.UserAuth, controller.MoviesList)
		User.GET("/shows", auth.UserAuth, controller.ShowsListByCinemas)
		User.GET("/show/seats", auth.UserAuth, controller.BookingLayout)
		User.POST("/show/booking", auth.UserAuth, controller.BookSeats)
		User.GET("/show/payment", controller.PaymentPage)
		User.POST("/booking/paymentSuceess", controller.PaymentValidation)
	}
}
