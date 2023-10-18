package controller

import (
	"context"
	"encoding/json"
	"net/http"
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

func UserProfile(c *gin.Context) {
	var user models.User
	username := c.GetString("username")

	result := config.DB.Select("first_name", "second_name", "email", "phone", "username", "status").First(&user, "username = ?", username)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Unable to get username"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": user})
}

func UserBookings(c *gin.Context) {
	var user models.User
	var bookings []models.Booking
	username := c.GetString("username")

	result := config.DB.Select("id", "first_name", "second_name", "email", "phone", "username").First(&user, "username = ?", username)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Unable to get username"})
		return
	}
	getBookings := config.DB.Where("user_id = ?", user.ID).Preload("Show").Preload("Show.Screen").Preload("Show.Screen.Cinemas").Preload("Show.Movie").Find(&bookings)
	if getBookings.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getBookings.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "bookings": bookings})
}

func UserPassChange(c *gin.Context) {
	var loginCred models.LoginCredentials
	if err := c.ShouldBindJSON(&loginCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	var user models.User
	result := config.DB.Where("username = ? OR email = ?", loginCred.Username, loginCred.Username).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Invalid Username/Email"})
		return
	}

	otp := GetOTP(user.FirstName, user.Email)
	//inserting the otp into reddis
	err := config.ReddisClient.Set(context.Background(), "forgetPassUserOTP"+user.Email, otp, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Error inserting otp in redis client"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Got to /change_password"})
}

func UserChangePassword(c *gin.Context) {
	var cred models.ChangePassword
	if err := c.ShouldBindJSON(&cred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	var user models.User
	result := config.DB.Where("username = ? OR email = ?", cred.Username, cred.Username).First(&user)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Invalid Username/Email"})
		return
	}

	if !verifyOTP("forgetPassUserOTP"+user.Email, cred.OTP, c) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid OTP"})
		return
	}

	pass, err := PassToHash(cred.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	result = config.DB.Model(&user).Update("password", string(pass))
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Password Changed Succesfully"})
}
