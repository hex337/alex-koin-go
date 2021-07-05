package command

import (
	"github.com/hex337/alex-koin-go/config"

	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func ProcessMessage(event *slackevents.AppMentionEvent) error {

	botID := fmt.Sprintf("<@%s> ", config.GetBotSlackID())

	name, err := parseCommandName(strings.TrimPrefix(event.Text, botID))
	if err != nil {
		return err
	}

	response, err := RunCommand(name, event)
	if err != nil {
		return err
	}

	err = replyWith(event.Channel, event.TimeStamp, response)
	if err != nil {
		return err
	}

	return nil
}

func parseCommandName(msg string) (string, error) {
	commands := map[string]string{
		"balance":  `^my balance`,
		"transfer": `^transfer to`,
	}
	for name, pattern := range commands {
		matched, err := regexp.MatchString(pattern, msg)
		if err != nil {
			return "", err
		}
		if matched {
			return name, nil
		}
	}
	return "", errors.New("no match")
}

func replyWith(channel string, msgTimestamp string, response string) error {
	botSecret := os.Getenv("SLACK_BOT_SECRET")
	var api = slack.New(botSecret)

	_, _, err := api.PostMessage(
		channel,
		slack.MsgOptionText(response, false),
		slack.MsgOptionTS(msgTimestamp), // reply in thread
	)
	return err
}
