package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Coin struct {
	gorm.Model
	Hash            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Origin          string
	MinedByUserID   uint
	UserID          uint
	CreatedByUserId uint

	Transactions []Transaction
}
