package controllers

import (
	"fmt"

	"github.com/first_project/database"
	"github.com/first_project/models"
	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context) {
	user, _ := c.Get("user")
	userid := user.(models.User).User_id

	var pid int

	err := database.DB.Model(&models.Payment{}).Select("payment_id").Where("user_id=?", userid).Scan(&pid).Error

	fmt.Println("payment id  : ", pid)

	if err != nil {
		c.HTML(400, "success.html", gin.H{
			"error": "Error in string conversion",
		})
		return
	}

	c.HTML(200, "success.html", gin.H{
		"paymentid": pid,
	})
}
