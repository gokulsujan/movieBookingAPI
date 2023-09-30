package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"theatreManagementApp/auth"
	"theatreManagementApp/config"
	"theatreManagementApp/models"
	"time"

	"github.com/gin-gonic/gin"
)

type otpCredentials struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type dateStr struct {
	DateStr string `json:"date"`
}

// User signup module
func UserSignUp(c *gin.Context) {
	inputField := models.User{}
	if err := c.ShouldBindJSON(&inputField); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Unable to bind json"})
		return
	}
	hash, err := PassToHash(inputField.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "False",
			"Error":  "Hashing password error",
		})
		return
	}
	inputField.Password = string(hash)
	var count int64

	config.DB.Model(&models.User{}).Where("Email = ?", inputField.Email).Count(&count)
	if count != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "This mail id is already registered with us"})
		return
	}
	config.DB.Model(&models.User{}).Where("Phone = ?", inputField.Phone).Count(&count)
	if count != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "This mobile number is already registered with us"})
		return
	}
	config.DB.Model(&models.User{}).Where("Username = ?", inputField.Username).Count(&count)
	if count != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Username already taken."})
		return
	}

	//generating otp and sending it to user
	Otp := GetOTP(inputField.FirstName, inputField.Email)

	jsonData, err := json.Marshal(inputField)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error encoding JSON"})
		return
	}

	//inserting the otp into reddis
	err = config.ReddisClient.Set(context.Background(), "signUpOTP"+inputField.Email, Otp, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting otp in redis client"})
		return
	}

	//inserting the data into reddis
	err = config.ReddisClient.Set(context.Background(), "userData"+inputField.Email, jsonData, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting user data in redis client"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "messsage": "Go to user/signup-verification"})
}

// user creation after email verification
func SignupVerification(c *gin.Context) {
	var otpCred otpCredentials
	if err := c.ShouldBindJSON(&otpCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err})
		return
	}

	if verifyOTP("signUpOTP"+otpCred.Email, otpCred.Otp, c) {
		var userData models.User
		superKey := "userData" + otpCred.Email
		jsonData, err := config.ReddisClient.Get(context.Background(), superKey).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error getting user data from redis client"})
			return
		}
		err = json.Unmarshal([]byte(jsonData), &userData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error binding reddis json data to user variable"})
			return
		} else {
			result := config.DB.Create(&userData)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
				return
			}
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Otp Verification success. User creation done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid OTP"})
	}
}

// user login module
func Userlogin(c *gin.Context) {
	var logincred models.LoginCredentials
	if err := c.ShouldBindJSON(&logincred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Unable to bind json data"})
	}
	var user models.User
	result := config.DB.First(&user, "username = ? OR email = ? OR phone = ?", logincred.Username, logincred.Username, logincred.Username)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid Credentials"})
		return
	}
	passMatch := HashToPass(user.Password, logincred.Password)
	if !passMatch {
		// Passwords do not match
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid password"})
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Not an active user. Contact customercare."})
		return
	}

	tokenString, err := auth.CreateToken(user.Username, "user")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "token": tokenString})
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
	cityName := c.Param("city")
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
	cityName := c.Param("city")
	var city models.City
	getCity := config.DB.Where("name ILIKE ?", "%"+cityName+"%").First(&city)
	if getCity.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": getCity.Error})
		return
	}
	var cinemas []models.Cinemas
	result := config.DB.Preload("City").Where("city_id = ?", city.ID).Find(&cinemas)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "CinemasList": cinemas})
}

func CinemasListOfMovies(c *gin.Context) {
	cityName := c.Param("city")
	movieId := c.Param("id")
	var city models.City
	getCity := config.DB.Where("name ILIKE ?", "%"+cityName+"%").First(&city)
	if getCity.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": getCity.Error})
		return
	}
	var cinemas []models.Cinemas
	result := config.DB.Preload("City").Distinct().Table("cinemas").Select("cinemas.*").Joins("JOIN screens ON cinemas.id = screens.cinemas_id").Joins("JOIN shows ON screens.id = shows.screen_id").Joins("JOIN movies ON shows.movie_id = movies.id").Where("cinemas.city_id = ? AND movies.id = ? AND shows.date >= ?", city.ID, movieId, time.Now().Format("2006-01-02")).Find(&cinemas)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "Error": result.Error})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "CinemasList": cinemas})
}

func ShowsListByCinemas(c *gin.Context) {
	movie_id := c.Param("id")
	cinemas_id := c.Param("cinemas")
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
	id := c.Param("id")
	var show models.Show

	result := config.DB.Preload("Screen").Preload("Screen.Cinemas").Preload("Screen.Cinemas.City").Preload("Screen.ScreenFormat").Preload("Movie").First(&show, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
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

func UserProfile(c *gin.Context) {
	var user models.User
	username := c.GetString("username")

	result := config.DB.Select("first_name", "second_name", "email", "phone", "username").First(&user, "username = ?", username)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Unable to get username"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": user})
}
