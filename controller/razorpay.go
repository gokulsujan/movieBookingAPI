package controller

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

func RazorpayOrderCreation(amt, book_id int) (string, error) {
	client := razorpay.NewClient(os.Getenv("razorPayApiId"), os.Getenv("razorPayApiPass"))

	data := map[string]interface{}{
		"amount":   amt * 100,
		"currency": "INR",
		"receipt":  strconv.Itoa(book_id),
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		return "", err
	}

	value := body["id"]
	str := value.(string)
	return str, nil
}

func PaymentPage(c *gin.Context) {
	id := c.Param("id")
	c.HTML(http.StatusAccepted, "payments.html", gin.H{"status": "true", "orderId": id})
}
