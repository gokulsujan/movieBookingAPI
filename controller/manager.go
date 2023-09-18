package controller

import (
	"net/http"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
)

func AddManager(c *gin.Context) {
	var manager models.Manager
	if err := c.ShouldBindJSON(&manager); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	hashedPass, err := PassToHash(manager.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status": "False",
			"Error":  "Hashing password error",
		})
		return
	}
	manager.Password = string(hashedPass)

	// checking the username already exists or not
	searchUsername := config.DB.First(&models.Manager{}, "Username=?", manager.Username)
	if searchUsername.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Username already exists for the manager role"})
		return
	}

	// checking the email already exists or not
	searchEmail := config.DB.First(&models.Manager{}, "Email=?", manager.Email)
	if searchEmail.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Email already exists for the manager role"})
		return
	}

	// checking already a manager assigned to the selected cinemas or not
	searchCinemas := config.DB.First(&models.Manager{}, "cinemas_id=?", manager.CinemasId)
	if searchCinemas.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Already another manager assigned to the selected cinemas"})
		return
	}

	result := config.DB.Create(&manager)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Manager created succesfully"})
}

func EditManager(c *gin.Context) {
	id := c.Param("id")
	var manager models.Manager
	if err := c.ShouldBindJSON(&manager); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": err.Error()})
		return
	}
	// checking the username already exists or not
	searchUsername := config.DB.Not("id = ?", id).First(&models.Manager{}, "Username=?", manager.Username)
	if searchUsername.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Username already exists for the manager role"})
		return
	}

	// checking the email already exists or not
	searchEmail := config.DB.Not("id = ?", id).First(&models.Manager{}, "Email=?", manager.Email)
	if searchEmail.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Email already exists for the manager role"})
		return
	}

	// checking already a manager assigned to the selected cinemas or not
	searchCinemas := config.DB.Not("id = ?", id).First(&models.Manager{}, "cinemas_id=?", manager.CinemasId)
	if searchCinemas.RowsAffected != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "false", "message": "Already another manager assigned to the selected cinemas"})
		return
	}
	result := config.DB.Model(&models.Manager{}).Where("id = ?", id).Updates(models.Manager{Name: manager.Name, Username: manager.Username, Email: manager.Email, CinemasId: manager.CinemasId, Status: manager.Status})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Manager details updated succefully"})

}

func DeleteManager(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Where("id = ?", id).Delete(&models.Manager{})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Manager deleted succesfully"})
}
