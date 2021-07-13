package command

import (
	"fmt"
)

type BalanceCommand struct{}

func (c *BalanceCommand) Run(msg string, event *CoinEvent) (string, error) {
	user := event.User

	results := fmt.Sprintf("Your current balance is %d :akc:", user.GetBalance())

	return results, nil
}
