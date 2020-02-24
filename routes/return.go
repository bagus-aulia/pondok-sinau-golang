package routes

import (
	"fmt"
	"strconv"

	"github.com/bagus-aulia/pondok-sinau-golang/config"
	"github.com/bagus-aulia/pondok-sinau-golang/models"
	"github.com/gin-gonic/gin"
)

//Return to handle book return
func Return(c *gin.Context) {
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

	if borrow.IsReturned {
		c.JSON(403, gin.H{
			"message": "Book has been returned. Cannot process return",
		})

		c.Abort()
		return
	}

	config.DB.Model(&borrow).First(&borrow, id).Updates(models.Transaction{
		IsReturned: true,
	})

	detailIDs := c.PostFormArray("detail_id[]")
	bookIDs := c.PostFormArray("book_id[]")
	fines := c.PostFormArray("fine[]")
	notes := c.PostFormArray("note[]")

	for i, detailID := range detailIDs {
		var detail models.Detail
		var book models.Book

		fine, _ := strconv.Atoi(fines[i])

		config.DB.Model(&detail).First(&detail, detailID).Updates(models.Detail{
			Fine: int32(fine),
			Note: notes[i],
		})

		if err := config.DB.Model(&book).First(&book, bookIDs[i]).Updates(map[string]interface{}{
			"Status": false,
		}).Error; err != nil {
			fmt.Println(err)
			return
		}
	}

	var newBorrow models.Transaction

	config.DB.Preload("Details").Preload("Details.Book").Preload("Admin").Preload("Member").First(&newBorrow, id)

	dataBorrow := genBorrowInfo(newBorrow)

	for _, newDetail := range newBorrow.Details {
		dataDetail := genDetailBorrowInfo(newDetail)

		dataBorrow.Detail = append(dataBorrow.Detail, dataDetail)
	}

	c.JSON(200, gin.H{
		"message": "Book has been returned",
		"data":    dataBorrow,
	})
}
