package routes

import (
	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

//ReturnGet to show borrow data by book code
func ReturnGet(c *gin.Context) {
	transCode := c.Param("code")
	var trans models.Transaction
	// var books models.Book

	// borrowLog := config.DB.Joins("JOIN detail_transactions ON detail_transactions.trans_id = transactions.id").Joins("JOIN admins ON admins.id = transactions.admin_id").Joins("JOIN members ON members.id = transactions.member_id").Joins("JOIN books ON books.id = detail_transactions.book_id AND books.code = ?", bookCode).First(&trans, "is_returned = ?", false)

	// borrowLog := config.DB.Preload("Detail_transactions").Preload("Admins").Preload("Members").Preload("Books", "code = ?", bookCode).First(&trans, "is_returned = ?", false)

	// borrowLog := config.DB.Preload("Details").Preload("Transactions", "code = ?", bookCode).Find(&books)

	borrowLog := config.DB.Preload("Details").Find(&trans, "code = ?", transCode)

	// var result models.All
	// borrowLog := config.DB.Preload("Details").Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)

	if borrowLog.RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Borrow log not found",
		})

		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   borrowLog,
	})
}

//ReturnUpdate to handle book return
func ReturnUpdate(c *gin.Context) {
	id := c.Param("id")
	var borrow models.Transaction

	if config.DB.First(&borrow, id).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Borrow log not found",
		})

		c.Abort()
		return
	}

	config.DB.Model(&borrow).First(&borrow, id).Updates(models.Transaction{
		IsReturned: true,
	})

	//update detail data

	c.JSON(200, gin.H{
		"message": "Book has been returned",
		"data":    borrow,
	})
}
