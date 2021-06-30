package command

import (
	"github.com/hex337/alex-koin-go/model"

	"fmt"
)

type BalanceCommand struct{}

func (c *BalanceCommand) Run(msg string) (string, error) {

	slackId := "W0122R46YBC"

	user, err := model.GetUserBySlackID(slackId)
	if err != nil {
		return "Internal error", err
	}

	results := fmt.Sprintf("Your current balance is %d koin", user.GetBalance())

	return results, nil
}
