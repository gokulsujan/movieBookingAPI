package controller

import (
	"fmt"
	"theatreManagementApp/config"
	"theatreManagementApp/models"
	"time"
)

func StartBookingMonitoring(bookingId uint) {
	for {
		var booking models.Booking
		getBooking := config.DB.First(&booking, bookingId)
		if getBooking.Error != nil {
			panic(getBooking.Error.Error())
		}
		if booking.Status == "success" {
			fmt.Printf("Booking %d is now successful. Exiting Goroutine.\n", booking.ID)
			return
		}
		if time.Since(booking.CreatedAt).Minutes() >= 10 {
			booking.Status = "cancelled"
			config.DB.Where("booking_id = ?", bookingId).Delete(&models.Seat{})
			config.DB.Save(&booking)
			fmt.Printf("Booking %d has been cancelled due to timeout\n", booking.ID)
			return // Exit the Goroutine
		}
		time.Sleep(1 * time.Minute)
	}
}
