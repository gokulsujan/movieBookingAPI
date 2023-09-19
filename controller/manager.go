package controller

import (
	"net/http"
	"theatreManagementApp/auth"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
)

func ManagerLogin(c *gin.Context) {
	var loginCred models.LoginCredentials
	var manager models.Manager
	if err := c.ShouldBindJSON(&loginCred); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	result := config.DB.Where("username = ? OR email = ?", loginCred.Username, loginCred.Username).First(&manager)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Invalid Credentials"})
		return
	}

	passMatch := HashToPass(manager.Password, loginCred.Password)
	if !passMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": "Incorrect Password"})
		return
	}

	token, err := auth.CreateToken(manager.Username, "manager")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "error": token})
}

func AddScreen(c *gin.Context) {
	var screen models.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	//checking the screen name already exists
	search := config.DB.First(&models.Screen{Name: screen.Name, CinemasId: screen.CinemasId})
	if search.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Screen Name already exists for the cinemas"})
		return
	}

	result := config.DB.Create(&screen)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Screen added to the cinemas"})
}

func EditScreen(c *gin.Context) {
	id := c.Param("id")
	var screen models.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	//checking the screen name already exists
	search := config.DB.Not("id = ?", id).First(&models.Screen{Name: screen.Name, CinemasId: screen.CinemasId})
	if search.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Screen Name already exists for the cinemas"})
		return
	}

	result := config.DB.Where("id = ?", id).Updates(&models.Screen{Name: screen.Name, CinemasId: screen.CinemasId, Rows: screen.Rows, Cols: screen.Cols, Premium: screen.Premium, Standard: screen.Standard, ScreenFormatId: screen.ScreenFormatId})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Screen Updated"})
}

func DeleteScreen(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&models.Screen{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Screen Deleted"})
}
