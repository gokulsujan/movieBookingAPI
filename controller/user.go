package controller

import (
	"context"
	"net/http"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type otpCredentials struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}

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

	//inserting the data into reddis
	config.ReddisClient.Set(context.Background(), "signUpData"+inputField.Email, inputField, 0)

	//inserting the otp into reddis
	config.ReddisClient.Set(context.Background(), "signUpOTP"+inputField.Email, Otp, 0)

	c.JSON(http.StatusAccepted, gin.H{"messsage": "Go to user/signup-verification"})
	// result := config.DB.Create(&inputField)
	// if result.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Error": result.Error})
	// 	return
	// }
	// c.JSON(http.StatusAccepted, gin.H{"success": "Data inserted into to the userdb"})
}

func SignupVerification(c *gin.Context) {
	var otpCred otpCredentials
	if err := c.ShouldBindJSON(&otpCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if verifyOTP("signUpOTP"+otpCred.Email, otpCred.Otp, c) {
		// inputField := models.User{}
		c.JSON(http.StatusAccepted, gin.H{"message": "Otp Verification done"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid OTP"})
	}
}
