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

	if config.DB.Preload("Details").Preload("Details.Book").Preload("Admin").Preload("Member").Find(&trans, "code = ?", transCode).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Borrow log not found",
		})

		c.Abort()
		return
	}

	var details []detailInfo

	for _, detailTrans := range trans.Details {
		detailStruct := detailInfo{
			ID:        detailTrans.ID,
			BookID:    detailTrans.BookID,
			BookCode:  detailTrans.Book.Code,
			BookTitle: detailTrans.Book.Title,
			BookCover: detailTrans.Book.Cover,
			Fine:      detailTrans.Fine,
			Note:      detailTrans.Note,
		}

		details = append(details, detailStruct)
	}

	c.JSON(200, gin.H{
		"id":             trans.ID,
		"code":           trans.Code,
		"adminID":        trans.AdminID,
		"adminUsername":  trans.Admin.Username,
		"memberID":       trans.MemberID,
		"memberUsername": trans.Member.Username,
		"borrowDate":     trans.BorrowDate,
		"returnDate":     trans.ReturnDate,
		"isReturned":     trans.IsReturned,
		"detail":         details,
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
