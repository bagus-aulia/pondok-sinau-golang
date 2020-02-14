package config

import (
	"os"

	"github.com/bagus-aulia/pondok-lentera/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //connect mysql
)

//DB is global variable to access database
var DB *gorm.DB

//InitDB used to connect to database function
func InitDB() {
	var err error

	DB, err = gorm.Open("mysql", os.Getenv("DB_USERNAME")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME")+"?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		panic("Failed to connect database")
	}

	DB.AutoMigrate(&models.Admin{})
	DB.AutoMigrate(&models.Book{})
	DB.AutoMigrate(&models.Member{})
	DB.AutoMigrate(&models.Transaction{}).AddForeignKey("admin_id", "admins(id)", "RESTRICT", "RESTRICT").AddForeignKey("member_id", "members(id)", "RESTRICT", "RESTRICT")
	DB.AutoMigrate(&models.DetailTransaction{}).AddForeignKey("trans_id", "transactions(id)", "RESTRICT", "RESTRICT").AddForeignKey("book_id", "books(id)", "RESTRICT", "RESTRICT")

	DB.Model(&models.Admin{}).Related(&models.Transaction{})
	DB.Model(&models.Member{}).Related(&models.Transaction{})
	DB.Model(&models.Transaction{}).Related(&models.DetailTransaction{})
	DB.Model(&models.Book{}).Related(&models.DetailTransaction{})
}
