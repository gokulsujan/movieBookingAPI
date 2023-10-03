package models

import "gorm.io/gorm"

type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePassword struct {
	OTP         string `json:"otp"`
	Username    string `json:"username"`
	NewPassword string `json:"newpassword"`
}

type Admin struct {
	gorm.Model
	Name     string `json:"name" gorm:"not null"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	IsSuper  bool   `json:"isSuper" gorm:"default:false"`
}

type Manager struct {
	gorm.Model
	Name      string  `json:"manager_name" gorm:"not null"`
	Username  string  `json:"username" gorm:"unique;not null"`
	Email     string  `json:"email", gorm:"unique;not null"`
	Password  string  `json:"password" gorm:"not null"`
	CinemasId uint    `json:"cinemas_id" gorm:"not null"`
	Cinemas   Cinemas `gorm:"not null"`
	Status    string  `json:"status" gorm:"default:active"`
}

type User struct {
	gorm.Model
	FirstName  string `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	SecondName string `json:"second_name" gorm:"not null" validate:"required,min=1,max=50"`
	Email      string `json:"email" gorm:"unique;not null"`
	Phone      string `json:"phone" gorm:"unique;not null"`
	Username   string `json:"username" gorm:"unique;not null"`
	Password   string `json:"password" gorm:"not null"`
	Status     string `json:"status" gorm:"default:active"`
}
