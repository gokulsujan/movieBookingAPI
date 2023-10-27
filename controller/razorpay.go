package controller

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorpayOrderCreation(amt, book_id int) (string, error) {
	client := razorpay.NewClient(os.Getenv("razorPayApiId"), os.Getenv("razorPayApiPass"))
	fmt.Println(book_id)
	OrderData := map[string]interface{}{
		"amount":   amt * 100,
		"currency": "INR",
		"receipt":  strconv.Itoa(book_id),
	}
	body, err := client.Order.Create(OrderData, nil)
	if err != nil {
		return "", err
	}

	value := body["id"]
	str := value.(string)
	fmt.Println("book id " + strconv.Itoa(book_id))
	return str, nil
}

func PaymentPage(c *gin.Context) {
	id := c.DefaultQuery("razorpay-order-id", "1")
	book_id := c.DefaultQuery("book-id", "0")
	if book_id == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "message": "Booking id is missing"})
		return
	}
	var booking models.Booking
	getBooking := config.DB.Preload("User").Preload("Show").Preload("Show.Movie").Preload("Show.Screen").Preload("Show.Screen.Cinemas").Preload("Show.Screen.ScreenFormat").First(&booking, book_id)
	if getBooking.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": getBooking.Error.Error()})
		return
	}
	var seats []models.Seat
	getSeats := config.DB.Where("booking_id = ?", book_id).Find(&seats)
	if getSeats.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": false, "error": getSeats.Error.Error()})
		return
	}
	if booking.Status == "success" {
		c.HTML(http.StatusAccepted, "bookingSuccess.html", gin.H{"status": "true", "orderId": id, "booking": booking, "seats": seats, "isBooked": true})
		return
	}
	c.HTML(http.StatusAccepted, "razorpay.html", gin.H{"status": "true", "orderId": id, "booking": booking, "seats": seats})
}

func PaymentValidation(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookingID, bookingIDExists := jsonData["notes"].(map[string]interface{})["booking_id"].(string)

	if !bookingIDExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID not found in the request"})
		return
	}
	result := config.DB.Model(&models.Booking{}).Where("id = ?", bookingID).Update("status", "success")
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	c.HTML(http.StatusAccepted, "bookingSuccess.html", gin.H{"status": "true", "message": "Booking succefull"})
}
