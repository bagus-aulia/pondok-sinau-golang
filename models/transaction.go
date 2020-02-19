package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Transaction models
type Transaction struct {
	gorm.Model
	Details    []Detail `gorm:"foreignkey:TransID"`
	Code       string   `gorm:"unique_index"`
	AdminID    uint
	MemberID   uint
	BorrowDate time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	ReturnDate time.Time
	IsReturned bool `gorm:"default:0"`
}
