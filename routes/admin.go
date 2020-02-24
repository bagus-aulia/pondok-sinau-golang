package routes

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

type adminInfo struct {
	ID       uint
	Username string
	FullName string
	Email    string
	Avatar   string
	IsAdmin  bool
}

//AdminIndex to show admin list
func AdminIndex(c *gin.Context) {
	admins := []models.Admin{}
	var returnAdmin []adminInfo

	config.DB.Find(&admins)

	for _, admin := range admins {
		dataAdmin := genAdminInfo(admin)

		returnAdmin = append(returnAdmin, dataAdmin)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   returnAdmin,
	})
}

//AdminProfile to show admin detail
func AdminProfile(c *gin.Context) {
	var admin models.Admin
	username := c.Param("username")

	if config.DB.First(&admin, "username = ?", username).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Account not found",
		})
		c.Abort()
		return
	}

	dataAdmin := genAdminInfo(admin)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   dataAdmin,
	})
}

//AdminUpdate to update admin data
func AdminUpdate(c *gin.Context) {
	//upload image not ready
	id := c.Param("id")
	var admin models.Admin
	var existAdmin models.Admin
	username := c.PostForm("username")
	email := c.PostForm("email")

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

	if !config.DB.Where("id <> ?", id).First(&existAdmin, "username = ?", username).RecordNotFound() {
		c.JSON(303, gin.H{
			"message": "Username already taken",
		})

		c.Abort()
		return
	}

	if !config.DB.Where("id <> ?", id).First(&existAdmin, "email = ?", email).RecordNotFound() {
		c.JSON(303, gin.H{
			"message": "Email already taken",
		})

		c.Abort()
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	newAvatarName := admin.Avatar

	if err == nil {
		dir, err := os.Getwd()
		fileLocation := filepath.Join(dir, "storage/admin", newAvatarName)
		err = os.Remove(fileLocation)

		filename := header.Filename
		extension := filepath.Ext(filename)
		random := rand.Intn(401)

		newAvatarName = admin.Username + "-" + strconv.Itoa(random) + extension
		fileLocation = filepath.Join(dir, "storage/admin", newAvatarName)

		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
				"note":    "error uploading image",
			})
		}

		defer targetFile.Close()

		if _, err := io.Copy(targetFile, file); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err,
				"note":    "error uploading image",
			})
		}
	}

	config.DB.Model(&admin).First(&admin, id).Updates(models.Admin{
		Username: c.PostForm("username"),
		FullName: c.PostForm("full_name"),
		Email:    c.PostForm("email"),
		Avatar:   newAvatarName,
	})

	dataAdmin := genAdminInfo(admin)

	c.JSON(200, gin.H{
		"message": "Account updated successfully",
		"data":    dataAdmin,
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

	dataAdmin := genAdminInfo(admin)

	c.JSON(200, gin.H{
		"message": "Account toggled successfully",
		"data":    dataAdmin,
	})
}

func genAdminInfo(admin models.Admin) adminInfo {
	dataAdmin := adminInfo{
		ID:       admin.ID,
		Username: admin.Username,
		FullName: admin.FullName,
		Email:    admin.Email,
		Avatar:   admin.Avatar,
		IsAdmin:  admin.IsAdmin,
	}

	return dataAdmin
}
