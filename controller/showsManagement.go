package controller

import (
	"fmt"
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
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Show added", "show-id": show.ID})
}

func GetRunnigMovies(c *gin.Context) {
	var shows []models.Show
	var dateStr dateStrShow
	if err := c.ShouldBindJSON(&dateStr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	result := config.DB.Preload("Movie").Preload("Screen").Preload("Screen.Cinemas").Preload("Screen.ScreenFormat").Preload("Screen.Cinemas.City").Where("DATE(date) = ? AND status = ?", dateStr.DateStr, "confirmed").Find(&shows)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "shows": shows, "date": dateStr.DateStr})
}

func ShowStatusChange(c *gin.Context) {
	type Status struct {
		Data string `json:"status"`
	}
	var status Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	fmt.Println(status.Data)
	if !(status.Data == "hold" || status.Data == "closed") {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "You can only hold the show or close the show. No other status update allowed"})
		return
	}
	id := c.DefaultQuery("show-id", "1")
	var show models.Show
	getShow := config.DB.First(&show, id)
	if getShow.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getShow.Error.Error()})
		return
	}
	if show.Status == "hold" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "This show is on hold. Admin should take action against the status of this show"})
		return
	}
	show.Status = status.Data
	result := config.DB.Save(&show)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "status updated to " + status.Data})
}

// Show Management by admin
func GetHoldedShows(c *gin.Context) {
	var shows []models.Show
	getShows := config.DB.Preload("Screen").Preload("Screen.Cinemas").Preload("Screen.ScreenFormat").Preload("Movie").Where("status = ?", "hold").Find(&shows)
	if getShows.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": getShows.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": true, "shows": shows})
}

func ChangeStatusShowByAdmin(c *gin.Context) {
	id := c.DefaultQuery("show-id", "0")
	var show models.Show
	getShow := config.DB.First(&show, id)
	if getShow.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": getShow.Error.Error()})
		return
	}

	type Status struct {
		Data string `json:"status"`
	}
	var status Status
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": err.Error()})
		return
	}

	show.Status = status.Data
	result := config.DB.Save(&show)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": true, "message": "Status updated"})
}

// refund initiater for cancelled bookings
func InitiateRefund(c *gin.Context) {
	id := c.DefaultQuery("show-id", "0")
	var show models.Show
	getShow := config.DB.Preload("Screen").Find(&show, id)
	if getShow.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": getShow.Error.Error()})
		return
	}
	var manager models.Manager
	username := c.GetString("username")
	getManager := config.DB.Where("username = ?", username).First(&manager)
	if getManager.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": getManager.Error.Error()})
		return
	}
	if show.Screen.CinemasId != manager.CinemasId {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Show is not assigned in your cinemas"})
		return
	}
	if show.Status != "cancell" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": false, "message": "Show is not yet cancelled"})
		return
	}
	var bookings []models.Booking
	getBookings := config.DB.Where("show_id = ? And Status = ?", show.ID, "success").Find(&bookings)
	if getBookings.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": getBookings.Error.Error()})
		return
	}
	for _, booking := range bookings {
		booking.Status = "cancell"

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
	}
	c.JSON(http.StatusAccepted, gin.H{"status": true, "message": "Refund initiated for all bookings"})
}
