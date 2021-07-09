package command

import (
	"github.com/hex337/alex-koin-go/model"

	"fmt"

	"github.com/slack-go/slack/slackevents"
)

type IdentityCommand struct{}

func (c *IdentityCommand) Run(msg string, event *slackevents.AppMentionEvent) (string, error) {
	slackId := event.User
	user, res, err := model.GetUserBySlackID(slackId)

	if !res {
		return fmt.Sprintf("User not found for %s.", slackId), nil
	} else if err != nil {
		return "Internal error", err
	}

	var message string

	role := user.Role()

	if role.Admin {
		message = "You are an Admin, you must keep the system going."
	} else if role.Lord {
		message = "You are a Lord of Koin, you must instill confidence in the system and exert control over the peasants."
	} else {
		message = "You are a peasant. Enjoy the Koin, and bask in it's glory."
	}

	return message, nil
}
