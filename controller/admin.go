package controller

import (
	"net/http"
	"theatreManagementApp/auth"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func AdminLogin(c *gin.Context) {
	var loginCred models.LoginCredentials
	if err := c.ShouldBindJSON(&loginCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Unable to bind json. The error: " + err.Error()})
		return
	}
	var adminUser models.Admin
	result := config.DB.First(&adminUser, "username = ?", loginCred.Username)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid credentials"})
		return
	}
	hashedPass := []byte(adminUser.Password)
	err := bcrypt.CompareHashAndPassword(hashedPass, []byte(loginCred.Password))
	if err != nil {
		// Passwords do not match
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Invalid password"})
		return
	}
	tokenString, err := auth.CreateToken(adminUser.Username, "admin")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "auth-token": tokenString})
}
