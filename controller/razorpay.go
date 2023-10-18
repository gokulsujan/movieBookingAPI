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
	c.HTML(http.StatusAccepted, "payments.html", gin.H{"status": "true", "orderId": id})
}

func PaymentValidation(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		// Handle JSON parsing error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the booking ID from the JSON data
	bookingID, bookingIDExists := jsonData["notes"].(map[string]interface{})["booking_id"].(string)

	if !bookingIDExists {
		// Handle the absence of booking ID
		c.JSON(http.StatusBadRequest, gin.H{"error": "Booking ID not found in the request"})
		return
	}
	result := config.DB.Model(&models.Booking{}).Where("id = ?", bookingID).Update("status", "success")
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	fmt.Println("Booking succefull")
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Booking succefull"})
}
