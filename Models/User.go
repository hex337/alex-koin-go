package Models

import (
	"github.com/hex337/alex-koin-go/Config"
)

func GetAllUsers(user *[]User) (err error) {
	if err = Config.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

