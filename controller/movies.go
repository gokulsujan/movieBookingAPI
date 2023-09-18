package controller

import (
	"net/http"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
)

func AddMovies(c *gin.Context) {
	var movie models.Movies
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	//checking the movie details already there or not
	searchMovie := config.DB.Where("name = ? OR description = ?", movie.Name, movie.Description).First(&movie)
	if searchMovie.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Movie already exists in the system"})
		return
	}
	result := config.DB.Create(&movie)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Movie added successfully"})
}
