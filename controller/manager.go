package controller

import (
	"fmt"
	"net/http"
	"strconv"
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

	token, err := auth.CreateToken(manager.Username, "manager", strconv.FormatUint(uint64(manager.CinemasId), 10))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "token": token})
}

func GetScreenList(c *gin.Context) {
	managerUsername := c.GetString("username")
	var manager models.Manager
	result := config.DB.Where("username = ?", managerUsername).First(&manager)
	if result.Error != nil {
		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
			return
		}
	}
	var screens []models.Screen

	result = config.DB.Where("cinemas_id = ?", manager.CinemasId).Find(&screens)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "screens": screens})
}

func AddScreen(c *gin.Context) {
	var screen models.Screen
	if err := c.ShouldBindJSON(&screen); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	fmt.Println(c.GetString("cinemas"))
	managerCinemas, err := strconv.Atoi(c.GetString("cinemas"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	if screen.CinemasId != uint(managerCinemas) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Unable to add screen to cinemas not assigned to you"})
		return
	}

	//checking the screen name already exists
	search := config.DB.Where("name = ? AND cinemas_id = ?", screen.Name, screen.CinemasId).First(&models.Screen{})
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
	search := config.DB.Not("id = ?", id).Where("name = ? AND cinemas_id = ?", screen.Name, screen.CinemasId).First(&models.Screen{})
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
	var manager models.Manager

	//getting manager details
	managerUsername := c.GetString("username")
	getManager := config.DB.Where("username = ?", managerUsername).First(&manager)
	if getManager.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Unautherised access"})
		return
	}
	var screen models.Screen
	//verifiying the manager is deleting the screen on his theathre
	verify := config.DB.Where("id = ?", id).First(&screen)
	if verify.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Already deleted"})
		return
	}

	if manager.CinemasId != screen.CinemasId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Unautherised access"})
		return
	}
	result := config.DB.Delete(&models.Screen{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Screen Deleted"})
}
