package controller

import (
	"encoding/json"
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
	return str, nil
}

func PaymentPage(c *gin.Context) {
	id := c.DefaultQuery("show-id", "1")
	c.HTML(http.StatusAccepted, "payments.html", gin.H{"status": "true", "orderId": id})
}

func PaymentValidation(c *gin.Context) {
	var callBackData models.RazorpayCallbackData
	if err := json.NewDecoder(c.Request.Body).Decode(&callBackData); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Extract custom booking ID from the callback data
	bookingID := callBackData.OrderData.BookingID
	result := config.DB.Model(&models.Booking{}).Where("id = ?", bookingID).Update("status", "success")
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}
	fmt.Println("Booking succefull")
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Booking succefull"})
}
