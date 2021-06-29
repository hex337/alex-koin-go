package model

import (
	"github.com/hex337/alex-koin-go/config"
)

func CreateUser(user *User) (err error) {
	if err = config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByID(user *User, id string) (err error) {
	if err = config.DB.Where("id = ?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserBySlackID(user *User, id string) (err error) {
	if err = config.DB.Where("slack_id = ?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

func GetAllUsers(user *[]User) (err error) {
	if err = config.DB.Find(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserBalance(user *User) (count int64) {
	return config.DB.Model(&user).Association("Coins").Count()
}
