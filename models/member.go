package models

import (
	"github.com/jinzhu/gorm"
)

//Member models
type Member struct {
	gorm.Model
	Transactions []Transaction
	Username     string `gorm:"unique_index"`
	FullName     string
	Phone        string
	Email        string `gorm:"unique_index"`
	Address      string `sql:"type:text;"`
	SocialID     string
	Provider     string
	Avatar       string
}
