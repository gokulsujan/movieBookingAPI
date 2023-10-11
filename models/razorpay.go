package models

type RazorpayCallbackData struct {
	OrderID   string `json:"order_id"`
	PaymentID string `json:"payment_id"`
	Amount    int64  `json:"amount"`
	Status    string `json:"status"`
	OrderData struct {
		BookingID string `json:"booking_id"`
	} `json:"notes"`
}
