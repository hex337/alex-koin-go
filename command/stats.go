package command

import (
	"fmt"

	"github.com/hex337/alex-koin-go/model"
)

type StatsCommand struct{}

// TODO: Add more stats - coins received, coins transfered, total coins
func (c *StatsCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
	user := event.User

	coinCount := model.CoinsCreatedThisWeekForUser(user)

	bot_response := BotResponse{
		Text: fmt.Sprintf("Stats for %s: ```Coins created this week: %d```", user.FirstName, coinCount),
	}

	return bot_response, nil
}
