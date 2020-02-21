package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Transaction models
type Transaction struct {
	gorm.Model
	Details    []Detail  `gorm:"foreignkey:TransID"`
	Admin      Admin     `gorm:"foreignkey:ID;association_foreignkey:AdminID"`
	Member     Member    `gorm:"foreignkey:ID;association_foreignkey:MemberID"`
	Code       string    `gorm:"unique_index"`
	AdminID    uint      `sql:"type:integer REFERENCES admins(id)"`
	MemberID   uint      `sql:"type:integer REFERENCES members(id)"`
	BorrowDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ReturnDate time.Time
	IsReturned bool `gorm:"default:0"`
}
