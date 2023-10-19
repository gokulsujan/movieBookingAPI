package controller

import (
	"net/http"
	"strconv"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
)

// City Management
func AddCity(c *gin.Context) {
	var city models.City

	if err := c.ShouldBindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	result := config.DB.First(&models.City{}, "Name  = ?", city.Name)
	if result.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "City name already exist"})
		return
	}

	result = config.DB.Create(&city)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": city.Name + " has been added in the city list", "city-id": city.ID})
}

func GetCityList(c *gin.Context) {
	var cities []models.City

	result := config.DB.Find(&cities)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "cityList": cities})
}

func EditCity(c *gin.Context) {
	id := c.Param("id")
	var city models.City
	result := config.DB.First(&city, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	var editedCity models.City
	if err := c.ShouldBindJSON(&editedCity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	city.Name = editedCity.Name
	result = config.DB.Save(&city)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": "Unable to update in the database"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "City name updated"})

}

func DeleteCity(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&models.City{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "City deleted successfully"})
}

//Cinemas Management

func AddCinemas(c *gin.Context) {
	var cinemas models.Cinemas
	if err := c.ShouldBindJSON(&cinemas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	search := config.DB.Where("name = ? AND city_id = ?", cinemas.Name, cinemas.CityId).First(&models.Cinemas{})
	if search.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Cinemas already exists in the city"})
		return
	}

	result := config.DB.Create(&cinemas)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Cinemas addedd successfully", "cinemas-id": cinemas.ID})
}

func EditCinemas(c *gin.Context) {
	id := c.Param("id")
	var cinemas models.Cinemas
	if err := c.ShouldBindJSON(&cinemas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	search := config.DB.Where("name = ? AND city_id = ?", cinemas.Name, cinemas.CityId).First(&models.Cinemas{})
	if search.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Cinemas already exists in the city"})
		return
	}
	result := config.DB.Model(&models.Cinemas{}).Where("id = ?", id).Updates(&models.Cinemas{Name: cinemas.Name, CityId: cinemas.CityId, Pincode: cinemas.Pincode})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Cinemas Updated"})
}

func DeleteCinemas(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&models.Cinemas{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Cinemas deleted succefully"})
}

func GetCinemasListFromCities(c *gin.Context) {
	city, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	var cinemas []models.Cinemas
	result := config.DB.Preload("City").Where("city_id = ?", uint(city)).Find(&cinemas)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "Cinemas List for the city ": cinemas})
}

// Adding ScreenFormat
func AddScreenFormat(c *gin.Context) {
	var format models.ScreenFormat
	if err := c.ShouldBindJSON(&format); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	//checking the format already exists
	search := config.DB.Where("screen_type = ? AND sound_system = ?", format.ScreenType, format.SoundSystem).First(&models.ScreenFormat{})
	if search.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Format already exists"})
		return
	}

	result := config.DB.Create(&format)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Screen format added", "screen-format-id": format.ID})

}

func ViewScreenFormat(c *gin.Context) {
	var screenFormats []models.ScreenFormat

	result := config.DB.Find(&screenFormats)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "screen-formats": screenFormats})
}
