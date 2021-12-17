package command

type IdentityCommand struct{}

func (c *IdentityCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
	user := event.User

	var message string

	role := user.Role()

	if role.Admin {
		message = "You are an Admin, you must keep the system going."
	} else if role.Lord {
		message = "You are a Lord of Koin, you must instill confidence in the system and exert control over the peasants."
	} else {
		message = "You are a peasant. Enjoy the Koin, and bask in it's glory."
	}

	return BotResponse{Text: message}, nil
}
