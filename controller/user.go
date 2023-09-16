package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"theatreManagementApp/config"
	"theatreManagementApp/models"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type otpCredentials struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

type loginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User signup module
func UserSignUp(c *gin.Context) {
	inputField := models.User{}
	if err := c.ShouldBindJSON(&inputField); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind json"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(inputField.Password), 10)
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
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This mail id is already registered with us"})
		return
	}
	config.DB.Model(&models.User{}).Where("Phone = ?", inputField.Phone).Count(&count)
	if count != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "This mobile number is already registered with us"})
		return
	}
	config.DB.Model(&models.User{}).Where("Username = ?", inputField.Username).Count(&count)
	if count != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Username already taken."})
		return
	}

	//generating otp and sending it to user
	Otp := GetOTP(inputField.FirstName, inputField.Email)

	jsonData, err := json.Marshal(inputField)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	//inserting the otp into reddis
	err = config.ReddisClient.Set(context.Background(), "signUpOTP"+inputField.Email, Otp, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error inserting otp in redis client"})
		return
	}

	//inserting the data into reddis
	err = config.ReddisClient.Set(context.Background(), "userData"+inputField.Email, jsonData, 1*time.Minute).Err()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error inserting user data in redis client"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"messsage": "Go to user/signup-verification"})
}

// user creation after email verification
func SignupVerification(c *gin.Context) {
	var otpCred otpCredentials
	if err := c.ShouldBindJSON(&otpCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if verifyOTP("signUpOTP"+otpCred.Email, otpCred.Otp, c) {
		var userData models.User
		superKey := "userData" + otpCred.Email
		jsonData, err := config.ReddisClient.Get(context.Background(), superKey).Result()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user data from redis client"})
			return
		}
		err = json.Unmarshal([]byte(jsonData), &userData)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error binding reddis json data to user variable"})
			return
		} else {
			result := config.DB.Create(&userData)
			if result.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"Error": result.Error})
				return
			}
		}

		c.JSON(http.StatusAccepted, gin.H{"message": "Otp Verification success. User creation done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
	}
}

// user login module
func Userlogin(c *gin.Context) {
	var logincred loginCredentials
	if err := c.ShouldBindJSON(&logincred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to bind json data"})
	}
	var user models.User
	result := config.DB.First(&user, "username = ? OR email = ? OR phone = ?", logincred.Username, logincred.Username, logincred.Username)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}
	hashedPass := []byte(user.Password)
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(logincred.Password))
	if err != nil {
		// Passwords do not match
		c.JSON(http.StatusUnauthorized, gin.H{"Message": "Invalid password"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Login done"})
}
