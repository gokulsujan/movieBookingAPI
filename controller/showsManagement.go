package controller

import (
	"net/http"
	"strconv"
	"theatreManagementApp/config"
	"theatreManagementApp/models"
	"time"

	"github.com/gin-gonic/gin"
)

type dateStrShow struct {
	DateStr string `json:"date"`
}

func AddShows(c *gin.Context) {
	var show models.Show
	if err := c.ShouldBindJSON(&show); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	//checking the manager is adding shows for the respective cinemas
	var screen models.Screen
	getScreenData := config.DB.First(&screen, show.ScreenId)
	if getScreenData.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getScreenData.Error.Error()})
		return
	}
	managerCinemas, err := strconv.Atoi(c.GetString("cinemas"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	if screen.CinemasId != uint(managerCinemas) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Unable to add screen to cinemas not assigned to you"})
		return
	}

	//checking is there any show is running for the selected time slot
	checkStartTime := show.Date.Add(-2 * time.Hour)
	checkEndTime := show.Date.Add(2 * time.Hour)
	searchResult := config.DB.Where("date BETWEEN ? AND ?", checkStartTime, checkEndTime).Where("screen_id = ?", show.ScreenId).First(&models.Show{})
	if searchResult.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Already show alloted in the screen"})
		return
	}

	result := config.DB.Create(&show)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Show added"})
}

func GetRunnigMovies(c *gin.Context) {
	var shows []models.Show
	var dateStr dateStrShow
	if err := c.ShouldBindJSON(&dateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	result := config.DB.Preload("Movie").Preload("Screen").Preload("Screen.Cinemas").Preload("Screen.ScreenFormat").Preload("Screen.Cinemas.City").Where("DATE(date) = ?", dateStr.DateStr).Find(&shows)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "shows": shows, "date": dateStr.DateStr})
}
