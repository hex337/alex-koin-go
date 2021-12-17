package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Nft struct {
	gorm.Model
	Hash            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	OwnedByUserId   uint
	CreatedByUserId uint
	MinBid          uint
	PricePaid       uint
	Name            string `gorm:"unique"`
	SourceURL       string
	DisplayURL      string
}
