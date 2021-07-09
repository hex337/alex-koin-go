package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string
	FirstName string
	LastName  string
	SlackID   string

	Coins        []Coin
	CoinsCreated []Coin `gorm:"foreignKey:CreatedByUserId"`

	TransactionsTo   []Transaction `gorm:"foreignKey:ToUserID"`
	TransactionsFrom []Transaction `gorm:"foreignKey:FromUserID"`
}

type UserRole struct {
	Admin bool
	Lord  bool
}
