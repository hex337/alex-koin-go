package command

import (
	"github.com/hex337/alex-koin-go/model"

	"fmt"

	"github.com/slack-go/slack/slackevents"
)

type BalanceCommand struct{}

func (c *BalanceCommand) Run(msg string, event *slackevents.AppMentionEvent) (string, error) {

	slackId := event.User

	user, res, err := model.GetUserBySlackID(slackId)

	if !res {
		return fmt.Sprintf("Your current balance is %d koin (user not found)", 0), nil
	} else if err != nil {
		return "Internal error", err
	}

	results := fmt.Sprintf("Your current balance is %d koin", user.GetBalance())

	return results, nil
}
