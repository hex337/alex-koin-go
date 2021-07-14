package model

import (
	"github.com/hex337/alex-koin-go/config"

	"errors"
	"os"
	"time"

	"github.com/slack-go/slack"
	"gorm.io/gorm"

	"github.com/davecgh/go-spew/spew"
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

func GetOrCreateUserBySlackID(slackId string) (*User, error) {
	user, res, err := GetUserBySlackID(slackId)

	if err != nil {
		return nil, err
	}

	if !res {
		botSecret := os.Getenv("SLACK_BOT_SECRET")
		var api = slack.New(botSecret)
		user, err := api.GetUserInfo(slackId)

		if err != nil {
			return nil, err
		}

		var dbUser User
		dbUser.FirstName = user.Profile.FirstName
		dbUser.LastName = user.Profile.LastName
		dbUser.SlackID = user.ID

		err = CreateUser(&dbUser)

		if err != nil {
			return nil, err
		}

		return &dbUser, nil
	}

	return user, nil
}

func (u *User) GetBalance() int64 {
	association := config.DB.Model(&u).Association("Coins")
	if association.Error != nil {
		return 0
	}

	return association.Count()
}

func (u *User) LastCoinInDays() (float64, error) {

	var lastCoin Coin

	whatAmI := config.DB.Model(&u).Order("created_at desc").Limit(1).Association("CoinsCreated").Find(&lastCoin)

	spew.Dump(whatAmI)

	// if err != nil {
	// return 0, err
	// }

	days := time.Now().Sub(lastCoin.CreatedAt).Hours() / 24
	spew.Dump(days)
	return days, nil
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
