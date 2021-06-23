package Model

import (
	"github.com/hex337/alex-koin-go/Config"
)

func CreateUser(user *User) (err error) {
	if err = Config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByID(user *User, id string) (err error) {
	if err = Config.DB.Preload("Coins").Where("id = ?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

func GetAllUsers(user *[]User) (err error) {
	if err = Config.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserBalance(user *User) (count int64) {
	return Config.DB.Model(&user).Association("Coins").Count()
}

