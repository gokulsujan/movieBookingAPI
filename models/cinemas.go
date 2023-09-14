package models

import (
	"gorm.io/gorm"
)

type City struct {
	gorm.Model
	Name string `json:"city" gorm:"not null"`
}

type Cinemas struct {
	gorm.Model
	Name    string `json:"cinemas_name" gorm:"not null"`
	CityId  uint   `json:"city_id" gorm:"not null"`
	City    City   `gorm:"ForeignKey:CityId"`
	Pincode string `json:"pincode" gorm:"not null"`
}

type Screen struct {
	gorm.Model
	Name           string       `json:"screen_name" gorm:"not null"`
	CinemasId      uint         `json:"cinemas_id" gorm:"not null"`
	Cinemas        Cinemas      `gorm:"ForeignKey:CinemasId"`
	Rows           int          `json:"rows" gorm:"not null"`
	Cols           int          `json:"cols" gorm:"not null"`
	Premium        int          `json:"prem" gorm:"not null"`
	Standard       int          `json:"std" gorm:"not null"`
	ScreenFormatId uint         `json:"screen_format" gorm:"not null"`
	ScreenFormat   ScreenFormat `gorm:"ForeignKey:ScreenFormatId"`
}

type ScreenFormat struct {
	gorm.Model
	ScreenType  string `json:"screen_type" gorm:"not null"`
	SoundSystem string `json:"sound_system" gorm:"not null"`
}
