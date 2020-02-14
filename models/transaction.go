package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//Transaction models
type Transaction struct {
	gorm.Model
	DetailTransactions []DetailTransaction
	Code               string `gorm:"unique_index"`
	AdminID            uint
	MemberID           uint
	BorrowDate         time.Time
	ReturnDate         time.Time
	IsReturned         bool `gorm:"default:0"`
}
