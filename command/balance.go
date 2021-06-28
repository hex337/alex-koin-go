package commands

type BalanceCommand struct{}

func (c *BalanceCommand) Run(msg string) (string, error) {

	return "Broke", nil
}
