package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	BookingId     uint    `json:"bookingid" gorm:"not null"`
	Booking       Booking `gorm:"ForeignKey:BookingId"`
	TotalAmount   uint    `json:"total_amount" gorm:"not null"`
	PaymentMethod string  `json:"payment_method" gorm:"not null"`
	Status        string  `json:"status" gorm:"not null"`
}

type Coupon struct {
	gorm.Model
	Code       string `json:"coupon_code" gorm:"not null"`
	Percent    uint   `json:"percent" gorm:"not null"`
	MaxPrice   uint   `json:"max_price" gorm:"not null"`
	ExpiryDate string `json:"expiry_date" gorm:"not null"`
	Status     string `json:"status" gorm:"not null"`
}
