package routes

import (
	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

//Home to handle homepage
func Home(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  200,
		"message": "this is homepage",
	})
}

//CheckUsername to check available username
func CheckUsername(c *gin.Context) {
	table := c.Param("role")
	username := c.PostForm("username")
	available := false
	status := "unavailable"

	if table == "admin" {
		var admin models.Admin
		if config.DB.First(&admin, "username = ?", username).RecordNotFound() {
			available = true
			status = "available"
		}
	} else {
		var member models.Member
		if config.DB.First(&member, "username = ?", username).RecordNotFound() {
			available = true
			status = "available"
		}
	}

	c.JSON(200, gin.H{
		"status":    200,
		"message":   "Username is " + status,
		"available": available,
	})
}
