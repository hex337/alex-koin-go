package model

import (
	"github.com/hex337/alex-koin-go/config"

	"errors"

	"gorm.io/gorm"
)

func CreateUser(user *User) (err error) {
	if err = config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByID(id int64) (*User, error) {
	var user User
	err := config.DB.Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, nil
	} else if err != nil {
		return &user, err
	}

	return &user, nil
}

func GetUserBySlackID(id string) (*User, error) {
	var user User
	err := config.DB.Where("slack_id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, nil
	} else if err != nil {
		return &user, err
	}

	return &user, nil
}

func (u *User) GetBalance() int64 {
	return config.DB.Model(&u).Association("Coins").Count()
}

func (u *User) IsEmpty() bool {
	return (&User{}) == u
}
