package model

import (
	"github.com/hex337/alex-koin-go/config"

	"github.com/nleeper/goment"
)

func CreateCoin(coin *Coin) (err error) {
	if err = config.DB.Create(&coin).Error; err != nil {
		return err
	}
	return nil
}

func CoinsCreatedThisWeekForUser(user *User) int64 {
	startOfWeek, _ := goment.New()
	startOfWeek.StartOf("week")
	dateFormat := startOfWeek.Format("YYYY-MM-DD")

	var count int64
	config.DB.Table("coins").Where("created_by_user_id = ? AND created_at > ?", user.ID, dateFormat).Count(&count)

	return count
}
