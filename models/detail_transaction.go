package models

import (
	"github.com/jinzhu/gorm"
)

//DetailTransaction models
type DetailTransaction struct {
	gorm.Model
	TransID uint
	BookID  uint
	Fine    int32
	Note    string `sql:"type:text;"`
}
