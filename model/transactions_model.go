package model

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Amount     int
	Memo       string
	FromUserID uint
	ToUserID   uint
	// CoinID     uint
}
