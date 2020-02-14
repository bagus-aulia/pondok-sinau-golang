package models

import (
	"github.com/jinzhu/gorm"
)

//User models
type User struct {
	gorm.Model
	Transactions []Transaction
	Username     string `gorm:"unique_index"`
	FullName     string
	Phone        string
	Email        string
	Address      string `sql:"type:text;"`
	SocialID     string
	Provider     string
	Avatar       string
	Role         bool `gorm:"default:0"`
}
