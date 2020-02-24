package routes

import (
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

type memberInfo struct {
	ID       uint
	Username string
	FullName string
	Phone    string
	Email    string
	Address  string
	Avatar   string
	Borrow   []borrowInfo
}

//MemberIndex to view member list
func MemberIndex(c *gin.Context) {
	members := []models.Member{}
	var returnMember []memberInfo

	config.DB.Find(&members)

	for _, member := range members {
		dataMember := genMemberInfo(member)

		returnMember = append(returnMember, dataMember)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   returnMember,
	})
}

//MemberProfile to view member detail
func MemberProfile(c *gin.Context) {
	var member models.Member
	username := c.Param("username")

	if config.DB.Preload("Transactions").Preload("Transactions.Details").Preload("Transactions.Details.Book").Preload("Transactions.Admin").First(&member, "username = ?", username).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Member not found",
		})
		c.Abort()
		return
	}

	dataMember := genMemberInfo(member)

	for _, borrow := range member.Transactions {
		dataBorrow := borrowInfo{
			ID:             borrow.ID,
			Code:           borrow.Code,
			AdminID:        borrow.AdminID,
			AdminUsername:  borrow.Admin.Username,
			MemberID:       borrow.MemberID,
			MemberUsername: member.Username,
			BorrowDate:     borrow.BorrowDate,
			ReturnDate:     borrow.ReturnDate,
			IsReturned:     borrow.IsReturned,
		}

		for _, detail := range borrow.Details {
			detailBorrow := detailInfo{
				ID:        detail.ID,
				BookID:    detail.BookID,
				BookCode:  detail.Book.Code,
				BookTitle: detail.Book.Title,
				Fine:      detail.Fine,
				Note:      detail.Note,
			}

			dataBorrow.Detail = append(dataBorrow.Detail, detailBorrow)
		}

		dataMember.Borrow = append(dataMember.Borrow, dataBorrow)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   dataMember,
	})
}

//MemberCreate to create new member
func MemberCreate(c *gin.Context) {
	var existMember models.Member
	username := slug.Make(c.PostForm("full_name"))
	email := c.PostForm("email")

	if !config.DB.First(&existMember, "username = ?", username).RecordNotFound() {
		c.JSON(303, gin.H{
			"message": "Username already taken",
		})

		c.Abort()
		return
	}

	if !config.DB.First(&existMember, "email = ?", email).RecordNotFound() {
		c.JSON(303, gin.H{
			"message": "Email already taken",
		})

		c.Abort()
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	newAvatarName := ""

	if err == nil {
		dir, err := os.Getwd()
		filename := header.Filename
		extension := filepath.Ext(filename)
		random := rand.Intn(401)

		newAvatarName = username + "-" + strconv.Itoa(random) + extension
		fileLocation := filepath.Join(dir, "storage/member", newAvatarName)

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

	member := models.Member{
		Username: slug.Make(c.PostForm("full_name")),
		FullName: c.PostForm("full_name"),
		Phone:    c.PostForm("phone"),
		Email:    c.PostForm("email"),
		Address:  c.PostForm("address"),
		Avatar:   newAvatarName,
	}

	config.DB.Create(&member)

	dataMember := genMemberInfo(member)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   dataMember,
	})
}

//MemberUpdate to update member data
func MemberUpdate(c *gin.Context) {
	id := c.Param("id")
	var member models.Member
	var existMember models.Member
	username := c.PostForm("username")
	email := c.PostForm("email")

	if config.DB.First(&member, id).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Member not found",
		})

		c.Abort()
		return
	}

	if !config.DB.Where("id <> ?", id).First(&existMember, "username = ?", username).RecordNotFound() {
		c.JSON(303, gin.H{
			"message": "Username already taken",
		})

		c.Abort()
		return
	}

	if !config.DB.Where("id <> ?", id).First(&existMember, "email = ?", email).RecordNotFound() {
		c.JSON(303, gin.H{
			"message": "Email already taken",
		})

		c.Abort()
		return
	}

	file, header, err := c.Request.FormFile("avatar")
	newAvatarName := member.Avatar

	if username == "" {
		username = member.Username
	}

	if err == nil {
		dir, err := os.Getwd()
		fileLocation := filepath.Join(dir, "storage/member", newAvatarName)
		err = os.Remove(fileLocation)

		filename := header.Filename
		extension := filepath.Ext(filename)
		random := rand.Intn(401)

		newAvatarName = username + "-" + strconv.Itoa(random) + extension
		fileLocation = filepath.Join(dir, "storage/member", newAvatarName)

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

	config.DB.Model(&member).First(&member, id).Updates(models.Member{
		Username: username,
		FullName: c.PostForm("full_name"),
		Phone:    c.PostForm("phone"),
		Email:    c.PostForm("email"),
		Address:  c.PostForm("address"),
		Avatar:   newAvatarName,
	})

	dataMember := genMemberInfo(member)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   dataMember,
	})
}

//MemberDelete to delete member data
func MemberDelete(c *gin.Context) {
	var member models.Member
	id := c.Param("id")

	memberFind := config.DB.First(&member, id)

	if memberFind.RecordNotFound() {
		c.JSON(404, gin.H{
			"message": "Member not found",
		})

		c.Abort()
		return
	}

	memberFind.Delete(&member)

	dataMember := genMemberInfo(member)

	c.JSON(200, gin.H{
		"message": "Member has been deleted",
		"data":    dataMember,
	})
}

func genMemberInfo(member models.Member) memberInfo {
	dataMember := memberInfo{
		ID:       member.ID,
		Username: member.Username,
		FullName: member.FullName,
		Phone:    member.Phone,
		Email:    member.Email,
		Address:  member.Address,
		Avatar:   member.Avatar,
	}

	return dataMember
}
