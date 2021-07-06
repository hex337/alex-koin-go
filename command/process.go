package command

import (
	"github.com/hex337/alex-koin-go/config"

	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func ProcessMessage(event *slackevents.AppMentionEvent) {

	botID := fmt.Sprintf("<@%s> ", config.GetBotSlackID())

	name, err := parseCommandName(strings.TrimPrefix(event.Text, botID))

	if err != nil {
		log.Printf("Could not parseCommandName : %v", err)
	}

	if name == "" {
		response := ":blob-wave: \n\nI am under construction and I am still learning how to handle koin.\n\nCheck out https://github.com/hex337/alex-koin-go"
		err := replyWith(event.Channel, event.TimeStamp, response)
		if err != nil {
			log.Printf("Could not replyWith : %v", err)
			return
		}
		return
	}

	response, err := RunCommand(name, event)
	if err != nil {
		log.Printf("Could not RunCommand : %v", err)
		return
	}

	err = replyWith(event.Channel, event.TimeStamp, response)
	if err != nil {
		log.Printf("Could not replyWith : %v", err)
		return
	}
	return
}

func parseCommandName(msg string) (string, error) {
	commands := map[string]string{
		"balance": `^my balance`,
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
	return "", nil
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
