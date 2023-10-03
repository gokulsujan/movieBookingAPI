package models

import "gorm.io/gorm"

type Booking struct {
	gorm.Model
	UserId   uint   `json:"uid" gorm:"not null"`
	User     User   `gorm:"ForeignKey:UserId"`
	ShowId   uint   `json:"showid" gorm:"not null"`
	Show     Show   `gorm:"ForeignKey:ShowId"`
	Seats    uint   `json:"seats" gorm:"not null"`
	Status   string `json:"status" gorm:"default:pending"`
	CouponId uint   `json:"coupon_code"`
	Coupon   Coupon `gorm:"ForeignKey:CouponId"`
}

type Seat struct {
	gorm.Model
	BookingId uint    `json:"booking_id" gorm:"not null"`
	Booking   Booking `gorm:"ForeignKey:BookingId"`
	SeatRow   int     `json:"seat_row" gorm:"not null"`
	SeatCol   int     `json:"seat_col" gorm:"not null"`
	Price     uint    `json:"Price" gorm:"not null"`
}
