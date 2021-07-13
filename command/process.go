package command

import (
	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"

	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type CoinEvent struct {
	User    *model.User
	Message string
}

func ProcessMessage(event *slackevents.AppMentionEvent) {

	botID := fmt.Sprintf("<@%s> ", config.GetBotSlackID())

	coinEvent, err := createCoinEvent(event)
	if err != nil {
		log.Printf("Error creating Coin Event : %v", err)
	}

	name, err := parseCommandName(strings.TrimPrefix(coinEvent.Message, botID))

	if err != nil {
		log.Printf("Could not parseCommandName : %v", err)
	}

	log.Printf("Command name : %v", name)

	if name == "" {
		response := ":blob-wave: \n\nI am under construction and I am still learning how to handle koin.\n\nCheck out https://github.com/hex337/alex-koin-go"
		err := replyWith(event.Channel, event.TimeStamp, response)
		if err != nil {
			log.Printf("Could not replyWith : %v", err)
			return
		}
		return
	}

	response, err := RunCommand(name, coinEvent)
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
		// Who says regexp are not readable
		"balance":     `(?i)^[[:space:]]*my[[:space:]]+balance.*`,
		"what_am_i":   `(?i)^[[:space:]]*what[[:space:]]+am[[:space:]]+i.*`,
		"create_coin": `(?i)^[[:space:]]*create[[:space:]]+koin.*`,
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

func createCoinEvent(event *slackevents.AppMentionEvent) (*CoinEvent, error) {
	botID := fmt.Sprintf("<@%s> ", config.GetBotSlackID())
	slackId := event.User

	var coinEvent CoinEvent

	trimmedMessage := strings.TrimPrefix(event.Text, botID)
	coinEvent.Message = trimmedMessage

	user, err := model.GetOrCreateUserBySlackID(slackId)

	if err != nil {
		log.Printf("[ERROR] Could not find or create a user for slack id %s", slackId)
		return &coinEvent, err
	}

	coinEvent.User = user

	log.Printf("Created CoinEvent: message - %s", coinEvent.Message)
	return &coinEvent, nil
}
