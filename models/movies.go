package models

import (
	"time"

	"gorm.io/gorm"
)

type Movies struct {
	gorm.Model
	Name           string `json:"movie_name" gorm:"not null"`
	Description    string `json:"description" gorm:"not null"`
	DurationMinute int    `json:"duration_minute" gorm:"not null"`
	ReleaseDate    string `json:"release_date" gorm:"not null"`
}

type Show struct {
	gorm.Model
	ScreenId uint      `json:"screenid" grom:"not null"`
	Screen   Screen    `gorm:"ForeignKey:ScreenId"`
	MovieId  uint      `json:"movieid" gorm:"not null"`
	Movie    Movies    `gorm:"ForeignKey:MovieId"`
	BaseFare uint      `json:"basefare" gorm:"not null"`
	Date     time.Time `json:"showdate" gorm:"not null"`
	Status   string    `json:"status" gorm:"not null"`
}
