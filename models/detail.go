package models

import (
	"github.com/jinzhu/gorm"
)

//Detail models
type Detail struct {
	gorm.Model
	Trans   Transaction `gorm:"foreignkey:ID;association_foreignkey:TransID"`
	Book    Book        `gorm:"foreignkey:ID;association_foreignkey:BookID"`
	TransID uint        `sql:"type:integer REFERENCES transactions(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	BookID  uint        `sql:"type:integer REFERENCES books(id)"`
	Fine    int32
	Note    string `sql:"type:text;"`
}
