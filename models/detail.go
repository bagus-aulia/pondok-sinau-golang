package models

import (
	"github.com/jinzhu/gorm"
)

//Detail models
type Detail struct {
	gorm.Model
	TransID uint `sql:"type:integer REFERENCES transactions(id) ON DELETE CASCADE ON UPDATE CASCADE"`
	BookID  uint `sql:"type:integer REFERENCES books(id)"`
	Fine    int32
	Note    string `sql:"type:text;"`
}
