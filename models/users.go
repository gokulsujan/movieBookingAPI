package models

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Username string `json:"username" gorm:"not null"`
	Password string `json:"password" gorm:"not null"`
}

type Manager struct {
	gorm.Model
	Name      string  `json:"manager_name" gorm:"not null"`
	Username  string  `json:"username" gorm:"not null"`
	Password  string  `json:"password" gorm:"not null"`
	CinemasId uint    `json:"cinemas_id" gorm:"not null"`
	Cinemas   Cinemas `gorm:"not null"`
}

type User struct {
	gorm.Model
	FirstName  string `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	SecondName string `json:"second_name" gorm:"not null" validate:"required,min=1,max=50"`
	Email      string `json:"email" gorm:"not null"`
	Phone      string `json:"phone" gorm:"not null"`
	Username   string `json:"username" gorm:"not null"`
	Password   string `json:"password" gorm:"not null"`
}
