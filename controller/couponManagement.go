package controller

import (
	"net/http"
	"theatreManagementApp/config"
	"theatreManagementApp/models"

	"github.com/gin-gonic/gin"
)

// Admin privileges
func GetCoupons(c *gin.Context) {
	var coupons []models.Coupon

	result := config.DB.Find(&coupons)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "coupons": coupons})
}

func AddCoupons(c *gin.Context) {
	var coupon models.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}

	//checking coupon code already exists or not
	getCoupon := config.DB.Where("code = ?", coupon.Code).First(&models.Coupon{})
	if getCoupon.RowsAffected != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "message": "Coupon code already exists"})
		return
	}

	result := config.DB.Create(&coupon)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Coupon created successfully", "coupon-id": coupon.ID})
}

func EditCoupon(c *gin.Context) {
	id := c.Param("id")
	var coupon models.Coupon
	if err := c.ShouldBindJSON(&coupon); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": err.Error()})
		return
	}
	result := config.DB.Where("id = ?", id).Updates(coupon)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Coupon updated successfully"})
}

func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")
	result := config.DB.Delete(&models.Coupon{}, id)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "false", "error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"status": "true", "message": "Coupon deleted successfully"})
}

//User privileges

func ValidateCoupon(code string) (uint, error) {
	var coupon models.Coupon
	result := config.DB.Where("code = ?", code).First(&coupon)
	if result.Error != nil {
		return 0, result.Error
	}
	return coupon.ID, nil
}

func CouponDiscountPrice(couponId uint, price int) int {
	var coupon models.Coupon
	config.DB.First(&coupon, couponId)
	discount := (price / 100) * int(coupon.Percent)
	if discount > int(coupon.MaxPrice) {
		return int(coupon.MaxPrice)
	}
	return discount
}
