package command

import (
	"fmt"

	"github.com/hex337/alex-koin-go/model"
)

type StatsCommand struct{}

// TODO: Add more stats - coins received, coins transfered, total coins
func (c *StatsCommand) Run(msg string, event *CoinEvent) (string, error) {
	user := event.User

	coinCount := model.CoinsCreatedThisWeekForUser(user)

	return fmt.Sprintf("Stats for %s: ```Coins created this week: %d```", user.FirstName, coinCount), nil
}
