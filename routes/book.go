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
)

//BookIndex to show book list
func BookIndex(c *gin.Context) {
	books := []models.Book{}

	config.DB.Find(&books)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   books,
	})
}

//BookDetail to show detail of book
func BookDetail(c *gin.Context) {
	var book models.Book
	code := c.Param("code")

	if config.DB.First(&book, "code = ?", code).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Book not found",
		})

		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   book,
	})
}

//BookCreate to handle create new book
func BookCreate(c *gin.Context) {
	//upload image cover not ready
	var lastBook models.Book
	var codeBook string

	if config.DB.Last(&lastBook).RecordNotFound() {
		codeBook = "LBK07110001"
	} else {
		lastCode := lastBook.Code
		runes := []rune(lastCode)
		first := string(runes[0:4])

		last := string(runes[4:])
		lastInt, _ := strconv.Atoi(last)
		lastInt = lastInt + 1

		codeBook = first + strconv.Itoa(lastInt)
	}

	file, header, err := c.Request.FormFile("cover")
	newCoverName := ""

	if err == nil {
		dir, err := os.Getwd()
		filename := header.Filename
		extension := filepath.Ext(filename)
		random := rand.Intn(401)

		newCoverName = codeBook + "-" + strconv.Itoa(random) + extension
		fileLocation := filepath.Join(dir, "storage/book", newCoverName)

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

	book := models.Book{
		Code:      codeBook,
		Title:     c.PostForm("title"),
		Publisher: c.PostForm("publisher"),
		Writer:    c.PostForm("writer"),
		Desc:      c.PostForm("desc"),
		Cover:     newCoverName,
	}

	config.DB.Create(&book)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   book,
	})
}

//BookUpdate to handle update book
func BookUpdate(c *gin.Context) {
	id := c.Param("id")
	var book models.Book

	if config.DB.First(&book, id).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Book not found",
		})

		c.Abort()
		return
	}

	file, header, err := c.Request.FormFile("cover")
	newCoverName := book.Cover

	if err == nil {
		dir, err := os.Getwd()
		fileLocation := filepath.Join(dir, "storage/book", newCoverName)
		err = os.Remove(fileLocation)

		filename := header.Filename
		extension := filepath.Ext(filename)
		random := rand.Intn(401)

		newCoverName = book.Code + "-" + strconv.Itoa(random) + extension
		fileLocation = filepath.Join(dir, "storage/book", newCoverName)

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

	config.DB.Model(&book).First(&book, id).Updates(models.Book{
		Title:     c.PostForm("title"),
		Publisher: c.PostForm("publisher"),
		Writer:    c.PostForm("writer"),
		Desc:      c.PostForm("desc"),
		Cover:     newCoverName,
	})

	c.JSON(200, gin.H{
		"status": 200,
		"data":   book,
	})
}

//BookDelete to handle delete book
func BookDelete(c *gin.Context) {
	var book models.Book
	id := c.Param("id")

	bookFind := config.DB.First(&book, id)

	if bookFind.RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Book not found",
		})

		c.Abort()
		return
	}

	bookFind.Delete(&book)

	c.JSON(200, gin.H{
		"message": "Book has been deleted",
		"data":    book,
	})
}
