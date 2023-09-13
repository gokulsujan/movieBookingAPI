package models

import "gorm.io/gorm"

type Movies struct {
	gorm.Model
	Name           string `json:"movie_name" gorm:"not null"`
	Description    string `json:"description" gorm:"not null"`
	DurationMinute int    `json:"duration_minute" gorm:"not null"`
	ReleaseDate    string `json:"release_date" gorm:"not null"`
}
