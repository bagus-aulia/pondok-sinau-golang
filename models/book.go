package models

import (
	"github.com/jinzhu/gorm"
)

//Book models
type Book struct {
	gorm.Model
	Code      string `gorm:"unique_index"`
	Title     string
	Publisher string
	Writer    string
	Desc      string
	Cover     string
	Status    bool `gorm:"default:0"`
}
