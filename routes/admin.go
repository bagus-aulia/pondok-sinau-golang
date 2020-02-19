package routes

import (
	"fmt"

	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

//AdminIndex to show admin list
func AdminIndex(c *gin.Context) {
	admins := []models.Admin{}

	config.DB.Find(&admins)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   admins,
	})
}

//AdminProfile to show admin detail
func AdminProfile(c *gin.Context) {
	var user models.Admin
	adminID := c.Param("id")

	admin := config.DB.First(&user, adminID)

	if admin.RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Account not found",
		})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   admin,
	})
}

//AdminUpdate to update admin data
func AdminUpdate(c *gin.Context) {
	id := c.Param("id")
	var admin models.Admin

	if config.DB.First(&admin, id).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Account not found",
		})

		c.Abort()
		return
	}

	if uint(c.MustGet("jwt_user_id").(float64)) != admin.ID && c.MustGet("jwt_is_admin") == false {
		c.JSON(403, gin.H{
			"status":  403,
			"message": "Forbidden access",
		})

		c.Abort()
		return
	}

	config.DB.Model(&admin).First(&admin, id).Updates(models.Admin{
		Username: c.PostForm("username"),
		FullName: c.PostForm("full_name"),
	})

	c.JSON(200, gin.H{
		"message": "Account updated successfully",
		"data":    admin,
	})
}

//AdminToggleRole used to set as admin or officer
func AdminToggleRole(c *gin.Context) {
	id := c.Param("id")
	var admin models.Admin

	if config.DB.First(&admin, id).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Account not found",
		})

		c.Abort()
		return
	}

	if err := config.DB.Model(&admin).First(&admin, id).Updates(map[string]interface{}{
		"IsAdmin": !admin.IsAdmin,
	}).Error; err != nil {
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"message": "Account toggled successfully",
		"data":    admin,
	})
}

//A
