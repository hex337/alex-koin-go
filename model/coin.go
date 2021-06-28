package model

import (
	"github.com/hex337/alex-koin-go/config"
)

func CreateCoin(coin *Coin) (err error) {
	if err = config.DB.Create(&coin).Error; err != nil {
		return err
	}
	return nil
}
