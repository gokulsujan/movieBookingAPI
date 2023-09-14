package controller

import (
	"net/http"
	"theatreManagementApp/config"
	"theatreManagementApp/models"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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
	Otp := GetOTP(inputField.FirstName, inputField.Email)
	session := sessions.Default(c)
	session.Set("expirationtime"+inputField.Email, time.Now().Add(time.Minute*1))
	session.Set("signUpOTP"+inputField.Email, Otp)
	session.Save()
	c.JSON(http.StatusAccepted, gin.H{"messsage": "Go to user/signup-otp-verification"})
	// result := config.DB.Create(&inputField)
	// if result.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Error": result.Error})
	// 	return
	// }
	// c.JSON(http.StatusAccepted, gin.H{"success": "Data inserted into to the userdb"})
}

func VerifyOTP() {

}
