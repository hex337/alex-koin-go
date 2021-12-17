package model

import (
	"log"

	"github.com/hex337/alex-koin-go/config"

	"errors"
	"os"

	"github.com/slack-go/slack"
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

func (u *User) GetCoin() (*Coin, error) {
	var coin Coin

	err := config.DB.Model(&u).Limit(1).Association("Coins").Find(&coin)
	if err != nil {
		return nil, err
	}

	return &coin, nil
}

func (u *User) GetNfts() ([]Nft, error) {
	var nfts []Nft

	err := config.DB.Table("nfts").Where("owned_by_user_id = ? AND price_paid > 0", u.ID).Find(&nfts).Error

	if err != nil {
		log.Printf("ERR error fetching nfts: %s", err)
		return nfts, err
	}

	return nfts, nil
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
