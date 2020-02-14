package models

import (
	"github.com/jinzhu/gorm"
)

//Admin models
type Admin struct {
	gorm.Model
	Transactions []Transaction
	Username     string `gorm:"unique_index"`
	FullName     string
	Email        string
	SocialID     string
	Provider     string
	Avatar       string
	IsAdmin      bool `gorm:"default:0"`
}
