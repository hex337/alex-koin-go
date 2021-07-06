package model

import (
	"github.com/hex337/alex-koin-go/config"

	"errors"
	"log"

	"gorm.io/gorm"
)

func CreateUser(user *User) (err error) {
	if err = config.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByID(id int64) (*User, bool, error) {
	var user User
	err := config.DB.Where("id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, false, nil
	} else if err != nil {
		return &user, false, err
	}

	return &user, true, nil
}

func GetUserBySlackID(id string) (*User, bool, error) {
	var user User
	err := config.DB.Where("slack_id = ?", id).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &user, false, nil
	} else if err != nil {
		return &user, false, err
	}

	return &user, true, nil
}

func (u *User) GetBalance() int64 {
	return config.DB.Model(&u).Association("Coins").Count()
}

func (u *User) IsEmpty() bool {
	empty := (&User{}) == u
	log.Printf("USER: %+v", u)
	log.Printf("USER: %+v", &User{})
	log.Printf("isEmpty: %t", empty)
	return empty
}
