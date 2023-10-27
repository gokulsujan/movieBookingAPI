package controller

import (
	"net/http"
	"strconv"
	"theatreManagementApp/config"
	"theatreManagementApp/models"
	"time"

	"github.com/gin-gonic/gin"
)

type BookingDetails struct {
	ShowBookingData models.Booking `json:"showDetails"`
	SelectedSeats   []models.Seat  `json:"selectedSeats"`
	CouponCode      string         `json:"coupon"`
}
type DateStruct struct {
	Datestr string `json:"date"`
}

func SelectCity(c *gin.Context) {
	var cities []models.City
	result := config.DB.Find(&cities)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "true", "Select City": cities})
}

func MoviesList(c *gin.Context) {
	cityName := c.DefaultQuery("city", "bengaluru")
	var dateStr dateStr
	if err := c.ShouldBindJSON(&dateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var city models.City
	getCity := config.DB.Where("name ILIKE ?", "%"+cityName+"%").First(&city)
	if getCity.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": getCity.Error})
		return
	}
	var movies []models.Movies
	result := config.DB.Distinct().Select("movies.name, movies.description, movies.duration_minute, movies.release_date").Joins("JOIN shows on movies.id = shows.movie_id").Joins("JOIN screens on screens.id = shows.screen_id").Joins("JOIN cinemas on cinemas.id = screens.cinemas_id").Where("cinemas.city_id = ? and DATE(shows.date) = ?", city.ID, dateStr.DateStr).Find(&movies)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "MoviesList": movies})
}

func CinemasList(c *gin.Context) {
	cityName := c.DefaultQuery("city", "bengaluru")
	movieId := c.DefaultQuery("movie-id", "")
	var city models.City
	getCity := config.DB.Where("name ILIKE ?", "%"+cityName+"%").First(&city)
	if getCity.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": getCity.Error})
		return
	}
	var cinemas []models.Cinemas
	if movieId == "" {
		result := config.DB.Preload("City").Where("city_id = ?", city.ID).Find(&cinemas)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
			return
		}
	} else {
		result := config.DB.Preload("City").Distinct().Table("cinemas").Select("cinemas.*").Joins("JOIN screens ON cinemas.id = screens.cinemas_id").Joins("JOIN shows ON screens.id = shows.screen_id").Joins("JOIN movies ON shows.movie_id = movies.id").Where("cinemas.city_id = ? AND movies.id = ? AND shows.date >= ?", city.ID, movieId, time.Now().Format("2006-01-02")).Find(&cinemas)
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
			return
		}
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "CinemasList": cinemas})
}

func CinemasListOfMovies(c *gin.Context) {
	cityName := c.DefaultQuery("city", "bengaluru")
	movieId := c.DefaultQuery("movie-id", "")
	var showDate DateStruct
	if err := c.ShouldBindJSON(&showDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if showDate.Datestr == "" {
		showDate.Datestr = time.Now().Format("2006-01-02")
	}
	var city models.City
	getCity := config.DB.Where("name ILIKE ?", "%"+cityName+"%").First(&city)
	if getCity.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": getCity.Error})
		return
	}
	var cinemas []models.Cinemas
	result := config.DB.Preload("City").Distinct().Table("cinemas").Select("cinemas.*").Joins("JOIN screens ON cinemas.id = screens.cinemas_id").Joins("JOIN shows ON screens.id = shows.screen_id").Joins("JOIN movies ON shows.movie_id = movies.id").Where("cinemas.city_id = ? AND movies.id = ? AND shows.date >= ?", city.ID, movieId, showDate.Datestr).Find(&cinemas)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "CinemasList": cinemas})
}

func ShowsListByCinemas(c *gin.Context) {
	movie_id := c.DefaultQuery("movies-id", "1")
	cinemas_id := c.DefaultQuery("cinemas-id", "1")
	var dateStr dateStr
	if err := c.ShouldBindJSON(&dateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var shows []models.Show
	result := config.DB.Preload("Screen").Preload("Screen.Cinemas").Preload("Screen.Cinemas.City").Preload("Screen.ScreenFormat").Preload("Movie").Table("shows").Joins("JOIN screens ON shows.screen_id = screens.id").Joins("JOIN cinemas ON screens.cinemas_id = cinemas.id").Joins("JOIN movies ON shows.movie_id = movies.id").Where("cinemas.id = ? AND movies.id = ? AND DATE(shows.date) = ?", cinemas_id, movie_id, dateStr.DateStr).Find(&shows)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "shows": shows})
}

// booking layout is the seating layout with booked and unbooked seats
func BookingLayout(c *gin.Context) {
	id := c.DefaultQuery("show-id", "1")
	var show models.Show

	result := config.DB.Preload("Screen").Preload("Screen.Cinemas").Preload("Screen.Cinemas.City").Preload("Screen.ScreenFormat").Preload("Movie").First(&show, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	if show.Status != "confirmed" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Bookings for this show is not available"})
		return
	}
	var bookedSeats []models.Seat
	result = config.DB.Table("seats").Joins("JOIN bookings ON seats.booking_id = bookings.id").Where("bookings.show_id = ?", show.ID).Find(&bookedSeats)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}

	seatAvailability := make([][]string, show.Screen.Rows)
	for i := range seatAvailability {
		seatAvailability[i] = make([]string, show.Screen.Cols)
	}

	for _, seat := range bookedSeats {
		seatAvailability[seat.SeatRow][seat.SeatCol] = "B"
	}
	seatLayout := make([][]string, show.Screen.Rows)
	for row := 0; row < show.Screen.Rows; row++ {
		seatLayout[row] = make([]string, show.Screen.Cols)
		for col := 0; col < show.Screen.Cols; col++ {
			if seatAvailability[row][col] == "" {
				seatLayout[row][col] = strconv.Itoa(col + 1) // Seat column number
			} else {
				seatLayout[row][col] = "B"
			}
		}
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "Seat": seatLayout})
}

func BookSeats(c *gin.Context) {
	var user models.User
	getUser := config.DB.Where("username = ?", c.GetString("username")).First(&user)
	if getUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "true", "error": getUser.Error.Error()})
		return
	}
	var booking BookingDetails
	if er := c.ShouldBindJSON(&booking); er != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "true", "error": er.Error()})
		return
	}
	var show models.Show
	config.DB.First(&show, booking.ShowBookingData.ShowId)
	if show.Status != "confirmed" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Bookings for this show is not available"})
		return
	}
	if booking.ShowBookingData.CouponId == 0 {
		booking.ShowBookingData.CouponId = 2
	}
	booking.ShowBookingData.UserId = user.ID

	bookResult := config.DB.Create(&booking.ShowBookingData)
	if bookResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": bookResult.Error.Error()})
		return
	}

	for i := range booking.SelectedSeats {
		checkSeat := config.DB.Where("seat_row = ? AND seat_col = ?", booking.SelectedSeats[i].SeatRow, booking.SelectedSeats[i].SeatCol).First(&models.Seat{})
		if checkSeat.RowsAffected != 0 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "false", "message": "Selected seat already on a booking"})
			return
		}
	}
	booking.ShowBookingData.Seats = uint(len(booking.SelectedSeats))
	price := 0
	for i := range booking.SelectedSeats {
		booking.SelectedSeats[i].BookingId = booking.ShowBookingData.ID
		price += int(booking.SelectedSeats[i].Price)
	}
	couponID, err := ValidateCoupon(booking.CouponCode)
	if err != nil {
		booking.ShowBookingData.CouponId = 2
	} else {
		booking.ShowBookingData.CouponId = couponID
	}

	discountPrice := CouponDiscountPrice(booking.ShowBookingData.CouponId, price)
	price -= discountPrice

	bookSeatResult := config.DB.Create(&booking.SelectedSeats)
	if bookSeatResult.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": bookSeatResult.Error.Error()})
		return
	}
	go StartBookingMonitoring(booking.ShowBookingData.ID)

	razorPayOrderId, err := RazorpayOrderCreation(price, int(booking.ShowBookingData.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Booking successfull, Payment pending. Wait until payment is successfull", "razorpayOrderId": razorPayOrderId, "booking-id": booking.ShowBookingData.ID})
}

func BookingCancellation(c *gin.Context) {
	id := c.DefaultQuery("book-id", "1")
	var booking models.Booking
	result := config.DB.Preload("Show").First(&booking, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	if booking.Show.Status != "confirmed" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "false", "message": "Unable to cancell unconfirmed show"})
		return
	}

	if booking.Status != "success" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "false", "message": "Unable to cancell unconfirmed bookings"})
		return
	}

	var user models.User
	getUser := config.DB.Where("username = ?", c.GetString("username")).First(&user)
	if getUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getUser.Error.Error()})
		return
	}
	if user.ID != booking.UserId {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "false", "message": "Unable to cancell other users booking"})
		return
	}

	//get seats
	var seats []models.Seat
	getSeats := config.DB.Where("booking_id = ?", booking.ID).Find(&seats)
	if getSeats.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getSeats.Error.Error()})
		return
	}
	cancelProcess := config.DB.Model(&models.Booking{}).Where("id = ?", booking.ID).Updates(&models.Booking{Status: "cancelled"})
	if cancelProcess.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": cancelProcess.Error.Error()})
		return
	}
	refundAmt := 0
	for i := range seats {
		refundAmt += int(seats[i].Price)
	}
	refundAmt -= CouponDiscountPrice(booking.CouponId, refundAmt)
	refundProcess := config.DB.Create(&models.Wallet{UserId: booking.UserId, Amt: refundAmt, Status: "success"})
	if refundProcess.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": refundProcess.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Booking Cancelled. Amoount after discount is refunded to the wallet."})
}
