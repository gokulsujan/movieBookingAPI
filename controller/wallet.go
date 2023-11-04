package controller

import (
	"fmt"
	"net/http"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
)

func GetBalance(c *gin.Context) {
	username := c.GetString("username")
	var user models.User
	getUser := config.DB.Where("username = ?", username).First(&user)
	if getUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getUser.Error.Error()})
		return
	}
	var transactions []models.Wallet
	fmt.Println(user.ID)
	result := config.DB.Where("user_id = ?", user.ID).Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	balance := 0
	for i := range transactions {
		balance += transactions[i].Amt

	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "balance": balance})
}

func PayWithWallet(c *gin.Context) {
	bookingId := c.DefaultQuery("book-id", "0")
	var booking models.Booking
	getBooking := config.DB.First(&booking, bookingId)
	if getBooking.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getBooking.Error.Error()})
		return
	}
	if booking.Status == "success" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "false", "message": "Payment already done for this booking"})
		return
	}
	var seats []models.Seat
	getSeats := config.DB.Where("booking_id = ?", booking.ID).Find(&seats)
	if getSeats.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "false", "error": getSeats.Error.Error()})
		return
	}
	amt := 0
	for i := range seats {
		amt += int(seats[i].Price)
	}

	//getting the wallet balance
	var transactions []models.Wallet
	result := config.DB.Where("user_id = ?", booking.UserId).Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	Walletbalance := 0
	for i := range transactions {
		Walletbalance += transactions[i].Amt
	}
	if Walletbalance < amt {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "false", "message": "Insufficient balance in the wallet"})
		return
	}
	paymentProcess := config.DB.Create(&models.Wallet{UserId: booking.UserId, Amt: 0 - amt, Status: "success"})
	if paymentProcess.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": paymentProcess.Error.Error()})
		return
	}
	updateBooking := config.DB.Model(&models.Booking{}).Where("id = ?", bookingId).Updates(&models.Booking{Status: "success"})
	if updateBooking.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": updateBooking.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Booking Successfull"})
}

func WalletTransactions(c *gin.Context) {
	username := c.GetString("username")
	var user models.User
	getUser := config.DB.Where("username = ?", username).First(&user)
	if getUser.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": getUser.Error.Error()})
		return
	}
	var transactions []models.Wallet
	result := config.DB.Where("user_id = ?", user.ID).Find(&transactions)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "transactions": transactions})
}
