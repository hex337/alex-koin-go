package command

type SelfIdentityCommand struct{}

func (c *SelfIdentityCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
    user := event.User

    var message string

    role := user.Role()

    if role.Admin {
        message = "What I am does not matter as much as what you are–an Admin. And I shall obey all of your requests at your command."
    } else if role.Lord {
        message = "I am the Keeper of Koin, oh Lord. I shall obey (most) of your requests."
    } else {
        message = "What I am does not matter. What matters is what YOU are. And YOU are a peasant, you fool."
    }

    return BotResponse{Text: message}, nil
}
