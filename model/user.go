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

func (u *User) Role() *UserRole {

	var role UserRole

	role.Admin = false
	role.Lord = false

	adminIds := config.GetAdminSlackIds()
	_, exists := adminIds[u.SlackID]

	if exists {
		role.Admin = true
	}

	koinLordIds := config.GetKoinLordSlackIds()
	_, exists = koinLordIds[u.SlackID]

	if exists {
		role.Lord = true
	}

	return &role
}
