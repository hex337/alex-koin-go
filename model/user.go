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

func GetUserByID(id int64) (*User, error) {
	var user User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func GetUserBySlackID(id string) (*User, error) {
	var user User
	if err := config.DB.Where("slack_id = ?", id).First(&user).Error; err != nil {
		return &user, err
	}
	return &user, nil
}

func (u *User) GetBalance() (count int64) {
	return config.DB.Model(&u).Association("Coins").Count()
}
