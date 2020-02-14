package routes

import (
	"strconv"
	"time"

	"github.com/bagus-train/gin-api-train/config"
	"github.com/bagus-train/gin-api-train/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

//GetIndex to handle index article page
func GetIndex(c *gin.Context) {
	articles := []models.Article{}
	config.DB.Find(&articles)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   articles,
	})
}

//GetArticlesByTag to show article by spesific tag
func GetArticlesByTag(c *gin.Context) {
	tag := c.Param("tag")
	articles := []models.Article{}

	config.DB.Where("tag LIKE ?", "%"+tag+"%").Find(&articles)

	c.JSON(200, gin.H{
		"data": articles,
	})
}

//GetArticle to handle readpage data
func GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	var article models.Article

	if config.DB.First(&article, "slug = ?", slug).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Article not found",
		})

		c.Abort()
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"data":   article,
	})
}

//CreateArticle to handle input data
func CreateArticle(c *gin.Context) {
	var findArticle models.Article
	slug := slug.Make(c.PostForm("title"))

	if !config.DB.First(&findArticle, "slug = ?", slug).RecordNotFound() {
		slug = slug + "-" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	article := models.Article{
		Title:   c.PostForm("title"),
		Content: c.PostForm("content"),
		Tag:     c.PostForm("tag"),
		Slug:    slug,
		UserID:  uint(c.MustGet("jwt_user_id").(float64)),
	}

	config.DB.Create(&article)

	c.JSON(200, gin.H{
		"status": 200,
		"data":   article,
	})
}

//UpdateArticle to handle update data
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if config.DB.First(&article, id).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Article not found",
		})

		c.Abort()
		return
	}

	if uint(c.MustGet("jwt_user_id").(float64)) != article.UserID {
		c.JSON(403, gin.H{
			"status":  403,
			"message": "Forbidden access",
		})

		c.Abort()
		return
	}

	config.DB.Model(&article).First(&article, id).Updates(models.Article{
		Title:   c.PostForm("title"),
		Tag:     c.PostForm("tag"),
		Content: c.PostForm("content"),
	})

	c.JSON(200, gin.H{
		"message": "berhasil update",
		"data":    article,
	})
}

//DeleteArticle to delete article
func DeleteArticle(c *gin.Context) {
	var article models.Article
	id := c.Param("id")

	if config.DB.First(&article, id).RecordNotFound() {
		c.JSON(404, gin.H{
			"status":  404,
			"message": "Article not found",
		})
	}

	config.DB.First(&article, id).Delete(&article)

	c.JSON(200, gin.H{
		"message": "Acticle has been deleted",
		"data":    article,
	})
}
