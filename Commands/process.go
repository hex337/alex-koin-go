package Commands

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
)

func ProcessMessage(channel string, msgTimestamp string, msg string) error {

	botID := "<@U025MPQCB9A> "

	name, err := parseCommandName(strings.TrimPrefix(msg, botID))
	if err != nil {
		return err
	}

	response, err := RunCommand(name)
	replyWith(channel, msgTimestamp, response)

	return err
}

func parseCommandName(msg string) (string, error) {
	commands := map[string]string{
		"balance":  `^my balance`,
		"transfer": `^transfer to`,
	}
	for name, pattern := range commands {
		matched, err := regexp.MatchString(pattern, msg)
		if err != nil {
			return "", nil
		}
		if matched {
			return name, nil
		}
	}
	return "", errors.New("no match")
}

func replyWith(channel string, msgTimestamp string, response string) (err error) {
	botSecret := os.Getenv("SLACK_BOT_SECRET")
	var api = slack.New(botSecret)

	api.PostMessage(
		channel,
		slack.MsgOptionText(response, false),
		slack.MsgOptionTS(msgTimestamp), // reply in thread
	)
	return err
}
