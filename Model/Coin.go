package Model

import (
	"github.com/hex337/alex-koin-go/Config"
)

func CreateCoin(coin *Coin) (err error) {
	if err = Config.DB.Create(&coin).Error; err != nil {
		return err
	}
	return nil
}
