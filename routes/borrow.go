package routes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

//BorrowIndex to show borrow list
func BorrowIndex(c *gin.Context) {
	borrow := []models.Transaction{}
	config.DB.Find(&borrow)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   borrow,
	})
}

//BorrowDetail to show borrow detail
func BorrowDetail(c *gin.Context) {
	code := c.Param("code")
	var borrow models.Transaction

	if config.DB.First(&borrow, "code = ?", code).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Borrow log not found",
		})

		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   borrow,
	})
}

//BorrowCreate to handle create borrow
func BorrowCreate(c *gin.Context) {
	var lastBorrow models.Transaction
	var codeBorrow string

	if config.DB.Last(&lastBorrow).RecordNotFound() {
		codeBorrow = "TPB13180001"
	} else {
		lastCode := lastBorrow.Code
		runes := []rune(lastCode)
		first := string(runes[0:3])

		last := string(runes[3:])
		lastInt, _ := strconv.Atoi(last)
		lastInt++

		codeBorrow = first + strconv.Itoa(lastInt)
	}

	memberID, _ := strconv.ParseUint(c.PostForm("member_id"), 10, 32)

	layoutDate := "2020-02-14"
	returnDt := c.PostForm("return_date")
	returnDate, _ := time.Parse(layoutDate, returnDt)
	// AdminID:    uint(memberID),

	borrow := models.Transaction{
		Code:       codeBorrow,
		AdminID:    uint(c.MustGet("jwt_user_id").(float64)),
		MemberID:   uint(memberID),
		ReturnDate: returnDate,
	}

	config.DB.Create(&borrow)

	details := c.PostFormArray("book_code")
	fmt.Println(details)
	fmt.Println("wewe")

	for _, bookCode := range details {
		var book models.Book

		if config.DB.First(&book, "code = ?", bookCode).RecordNotFound() {
			fmt.Println("error")
			continue
		}

		fmt.Println(book.ID)
	}

	//detail form not complete
}

//BorrowUpdate will handle update borrow log
func BorrowUpdate(c *gin.Context) {
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

	memberID, _ := strconv.ParseUint(c.PostForm("member_id"), 10, 32)

	layoutDate := "2020-02-14 18:55:41"
	returnDt := c.PostForm("return_date")
	returnDate, _ := time.Parse(layoutDate, returnDt)

	config.DB.Model(&borrow).First(&borrow, id).Updates(models.Transaction{
		MemberID:   uint(memberID),
		ReturnDate: returnDate,
	})

	//update detail data

	c.JSON(200, gin.H{
		"message": "Borrow log has been updated",
		"data":    borrow,
	})
}

//BorrowDelete to handle deleting borrow log
func BorrowDelete(c *gin.Context) {
	id := c.Param("id")
	var borrow models.Transaction
	borrowLog := config.DB.First(&borrow, id)

	if borrowLog.RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Borrow log not found",
		})

		c.Abort()
		return
	}

	borrowLog.Delete(&borrow)

	//delete detail data

	c.JSON(200, gin.H{
		"message": "Borrow log has been deleted",
		"data":    borrow,
	})
}
