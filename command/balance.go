package command

import (
	"github.com/hex337/alex-koin-go/model"

	"fmt"
)

type BalanceCommand struct{}

func (c *BalanceCommand) Run(msg string) (string, error) {

	slackId := "W0122R46YBC"

	var user model.User
	err := model.GetUserBySlackID(&user, slackId)
	if err != nil {
		return "Internal error", err
	}

	balance := model.GetUserBalance(&user)
	results := fmt.Sprintf("Your current balance is %d koin", balance)

	return results, nil
}
