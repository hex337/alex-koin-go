package command

import (
	"github.com/hex337/alex-koin-go/model"

	"fmt"

	"github.com/slack-go/slack/slackevents"
)

type BalanceCommand struct{}

func (c *BalanceCommand) Run(msg string, event *slackevents.AppMentionEvent) (string, error) {

	slackId := event.User

	user, err := model.GetUserBySlackID(slackId)

	if user.IsEmpty() {
		return fmt.Sprintf("Your current balance is %d koin", 0), nil
	} else if err != nil {
		return "Internal error", err
	}

	results := fmt.Sprintf("Your current balance is %d koin", user.GetBalance())

	return results, nil
}
