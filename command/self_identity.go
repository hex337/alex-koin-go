package command

type SelfIdentityCommand struct{}

func (c *SelfIdentityCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
    user := event.User

    var message string

    role := user.Role()

    if role.Admin {
        message = "I am your sovereign ruler, Alex Koin, guiding and overseeing the administration to ensure the kingdom's efficient operation."
    } else if role.Lord {
        message = "I am your esteemed ruler, Alex Koin, embodying supreme authority and wisdom, providing guidance and leadership to uphold the kingdom's prosperity."
    } else {
        message = "I am your benevolent ruler, Alex Koin. I bring joy to your life. Never forget to honor your lords."
    }

    return BotResponse{Text: message}, nil
}
