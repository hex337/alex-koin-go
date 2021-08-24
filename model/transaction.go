package model

import (
	"github.com/hex337/alex-koin-go/config"
)

func CreateTransaction(transaction *Transaction) (err error) {
	if err = config.DB.Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}
