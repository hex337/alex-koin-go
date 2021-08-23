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
	User           *model.User
	Channel        string
	Message        string
	ReplyTimeStamp string
}

func ProcessMessageEvent(event *slackevents.MessageEvent) {
	botID := fmt.Sprintf("<@%s> ", config.GetBotSlackID())
	noPrefixChannels := config.GetNoPrefixChannelIds()

	_, exists := noPrefixChannels[event.Channel]

	if !exists {
		// we need to prefix with bot id then
		if !strings.HasPrefix(event.Text, botID) {
			return
		}
	}

	coinEvent, err := createCoinEventFromMessageEvent(event)
	if err != nil {
		log.Printf("Error creating Coin Event : %v", err)
	}

	executeCoinEvent(coinEvent)
	return
}

func ProcessAppMentionEvent(event *slackevents.AppMentionEvent) {

	// coinEvent, err := createCoinEventFromAppMention(event)
	// if err != nil {
	// 	log.Printf("Error creating Coin Event : %v", err)
	// }

	// executeCoinEvent(coinEvent)
	return
}

func executeCoinEvent(coinEvent *CoinEvent) {
	name, err := parseCommandName(coinEvent.Message)
	if err != nil {
		log.Printf("Could not parseCommandName : %v", err)
	}

	log.Printf("Command name : %v", name)

	if name == "" {
		log.Printf("Command name empty, no actions to be had.")
		return
	}

	response, err := RunCommand(name, coinEvent)
	if err != nil {
		log.Printf("Could not RunCommand : %v", err)
		return
	}

	err = replyWith(coinEvent.Channel, coinEvent.ReplyTimeStamp, response)
	if err != nil {
		log.Printf("Could not replyWith : %v", err)
		return
	}
	return
}

func parseCommandName(msg string) (string, error) {
	log.Printf("Attempting to match message: '%s'", msg)

	commands := map[string]string{
		// Who says regexp are not readable
		"balance":      `(?i)^[[:space:]]*my[[:space:]]+balance.*`,
		"what_am_i":    `(?i)^[[:space:]]*what[[:space:]]+am[[:space:]]+i.*`,
		"create_coin":  `(?i)^[[:space:]]*create[[:space:]]+koin.*`,
		"stats":        `(?i)^[[:space:]]*stats$`,
		"destroy_koin": `(?i)^[[:space:]]*destroy[[:space:]]+koin.*`,
	}
	for name, pattern := range commands {
		matched, err := regexp.MatchString(pattern, msg)
		if err != nil {
			log.Printf("Error with regex: %v", err)
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

func createCoinEventFromMessageEvent(event *slackevents.MessageEvent) (*CoinEvent, error) {
	return createCoinEvent(event.User, event.Channel, event.Text, event.ThreadTimeStamp, event.TimeStamp)
}

func createCoinEventFromAppMention(event *slackevents.AppMentionEvent) (*CoinEvent, error) {
	return createCoinEvent(event.User, event.Channel, event.Text, event.ThreadTimeStamp, event.TimeStamp)
}

func createCoinEvent(slackId string, channel string, message string, threadTimeStamp string, timeStamp string) (*CoinEvent, error) {
	botID := fmt.Sprintf("<@%s> ", config.GetBotSlackID())

	var coinEvent CoinEvent
	var replyTimeStamp string

	coinEvent.Channel = channel

	if threadTimeStamp != "" {
		// reply in the thread
		replyTimeStamp = threadTimeStamp
	} else {
		// reply in a new thread on the original comment message
		replyTimeStamp = timeStamp
	}

	coinEvent.ReplyTimeStamp = replyTimeStamp

	trimmedMessage := strings.TrimPrefix(message, botID)
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
