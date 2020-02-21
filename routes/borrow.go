package routes

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

type borrowInfo struct {
	ID             uint
	Code           string
	AdminID        uint
	AdminUsername  string
	MemberID       uint
	MemberUsername string
	BorrowDate     time.Time
	ReturnDate     time.Time
	IsReturned     bool
	Detail         []detailInfo
}

type detailInfo struct {
	ID        uint
	BookID    uint
	BookCode  string
	BookTitle string
	BookCover string
	Fine      int32
	Note      string
}

//BorrowIndex to show borrow list
func BorrowIndex(c *gin.Context) {
	borrows := []models.Transaction{}
	var dataBorrows []borrowInfo

	config.DB.Preload("Details").Preload("Details.Book").Preload("Admin").Preload("Member").Find(&borrows)

	for _, borrow := range borrows {
		dataBorrow := genBorrowInfo(borrow)

		for _, detail := range borrow.Details {
			dataDetail := genDetailBorrowInfo(detail)

			dataBorrow.Detail = append(dataBorrow.Detail, dataDetail)
		}

		dataBorrows = append(dataBorrows, dataBorrow)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   dataBorrows,
	})
}

//BorrowDetail to show borrow detail
func BorrowDetail(c *gin.Context) {
	code := c.Param("code")
	var borrow models.Transaction

	if config.DB.Preload("Details").Preload("Details.Book").Preload("Admin").Preload("Member").First(&borrow, "code = ?", code).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Borrow log not found",
		})

		c.Abort()
		return
	}

	dataBorrow := genBorrowInfo(borrow)

	for _, detail := range borrow.Details {
		dataDetail := genDetailBorrowInfo(detail)

		dataBorrow.Detail = append(dataBorrow.Detail, dataDetail)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   dataBorrow,
	})
}

//BorrowCreate to handle create borrow
func BorrowCreate(c *gin.Context) {
	//need improvement
	//check return date saved or not
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

	layoutDate := "2020-02-14 18:52:17"
	returnDt := c.PostForm("return_date")
	returnDate, _ := time.Parse(layoutDate, returnDt)

	borrow := models.Transaction{
		Code:       codeBorrow,
		AdminID:    uint(c.MustGet("jwt_user_id").(float64)),
		MemberID:   uint(memberID),
		ReturnDate: returnDate,
	}

	config.DB.Create(&borrow)

	details := c.PostFormArray("book_code[]")

	for _, bookCode := range details {
		var book models.Book

		if config.DB.First(&book, "code = ?", bookCode).RecordNotFound() {
			fmt.Println("error")
			continue
		}

		detail := models.Detail{
			TransID: borrow.ID,
			BookID:  book.ID,
		}

		config.DB.Create(&detail)
	}

	var newBorrow models.Transaction
	config.DB.Preload("Details").Preload("Details.Book").Preload("Admin").Preload("Member").First(&newBorrow, borrow.ID)

	dataBorrow := genBorrowInfo(newBorrow)

	for _, newDetail := range newBorrow.Details {
		dataDetail := genDetailBorrowInfo(newDetail)

		dataBorrow.Detail = append(dataBorrow.Detail, dataDetail)
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   dataBorrow,
	})
}

//BorrowUpdate will handle update borrow log
func BorrowUpdate(c *gin.Context) {
	//need improvement
	//check update detail
	//check return date saved or not
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

	details := c.PostFormArray("book_code[]")
	detailIDs := c.PostFormArray("detail_id[]")

	for i, bookCode := range details {
		var book models.Book
		var detail models.Detail

		if config.DB.First(&book, "code = ?", bookCode).RecordNotFound() {
			fmt.Println("error")
			continue
		}

		config.DB.Model(&detail).First(&detail, detailIDs[i]).Updates(models.Detail{
			BookID: book.ID,
		})
	}

	var newBorrow models.Transaction

	config.DB.Preload("Details").Preload("Details.Book").Preload("Admin").Preload("Member").First(&newBorrow, id)

	dataBorrow := genBorrowInfo(newBorrow)

	for _, newDetail := range newBorrow.Details {
		dataDetail := genDetailBorrowInfo(newDetail)

		dataBorrow.Detail = append(dataBorrow.Detail, dataDetail)
	}

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

func genBorrowInfo(borrow models.Transaction) borrowInfo {
	dataBorrow := borrowInfo{
		ID:             borrow.ID,
		Code:           borrow.Code,
		AdminID:        borrow.AdminID,
		AdminUsername:  borrow.Admin.Username,
		MemberID:       borrow.MemberID,
		MemberUsername: borrow.Member.Username,
		BorrowDate:     borrow.BorrowDate,
		ReturnDate:     borrow.ReturnDate,
		IsReturned:     borrow.IsReturned,
	}

	return dataBorrow
}

func genDetailBorrowInfo(detail models.Detail) detailInfo {
	dataDetail := detailInfo{
		ID:        detail.ID,
		BookID:    detail.BookID,
		BookCode:  detail.Book.Code,
		BookTitle: detail.Book.Title,
		BookCover: detail.Book.Cover,
		Fine:      detail.Fine,
		Note:      detail.Note,
	}

	return dataDetail
}
