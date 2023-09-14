package config

import (
	"os"
	"theatreManagementApp/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.Admin{})
	DB.AutoMigrate(&models.Manager{})
	DB.AutoMigrate(&models.User{})

	DB.AutoMigrate(&models.City{})
	DB.AutoMigrate(&models.Cinemas{})
	DB.AutoMigrate(&models.Screen{})
	DB.AutoMigrate(&models.ScreenFormat{})

	DB.AutoMigrate(&models.Movies{})
	DB.AutoMigrate(&models.Show{})

	DB.AutoMigrate(&models.Booking{})
	DB.AutoMigrate(&models.Seat{})

	DB.AutoMigrate(&models.Payment{})
	DB.AutoMigrate(&models.Coupon{})
}
