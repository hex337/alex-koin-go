package model

import (
	"github.com/hex337/alex-koin-go/config"
	// "github.com/nleeper/goment"
)

func CreateNft(nft *Nft) (err error) {
	if err = config.DB.Create(&nft).Error; err != nil {
		return err
	}
	return nil
}

func (n *Nft) DestroyNft() (err error) {
	if err = config.DB.Delete(n).Error; err != nil {
		return err
	}
	return nil
}
