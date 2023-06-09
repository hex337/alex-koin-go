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

type BotResponse struct {
	Pretext        string
	Text           string
	Fields         []slack.AttachmentField
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

	bot_response, err := RunCommand(name, coinEvent)
	if err != nil {
		log.Printf("Could not RunCommand : %v", err)
		return
	}

	err = replyWith(coinEvent.Channel, coinEvent.ReplyTimeStamp, bot_response.Text)
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
		"all_nfts":      `(?i)^[[:space:]]*all[[:space:]]+nfts$`,
		"balance":       `(?i)^[[:space:]]*my[[:space:]]+balance.*`,
		"create_coin":   `(?i)^[[:space:]]*create[[:space:]]+koin.*`,
		"create_nft":    `(?i)^[[:space:]]*create[[:space:]]+nft.*`,
		"destroy_coin":  `(?i)^[[:space:]]*destroy[[:space:]]+koin.*`,
		"my_nfts":       `(?i)^[[:space:]]*my[[:space:]]+nfts$`,
		"stats":         `(?i)^[[:space:]]*stats$`,
		"transfer_coin": `(?i)^[[:space:]]*transfer[[:space:]]+[0-9]+[[:space:]].*`,
		"transfer_nft":  `(?i)^[[:space:]]*transfer[[:space:]]+nft+[[:space:]].*`,
		"what_am_i":     `(?i)^[[:space:]]*what[[:space:]]+am[[:space:]]+i.*`,
		"what_are_you":  `(?i)^[[:space:]]*what[[:space:]]+are[[:space:]]+you.*`,
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

	nbsp := "\u00A0"
	trimmedMessage := strings.TrimPrefix(message, botID)
	// If you copy and paste in slack, it adds non-breaking spaces around @ mentions
	coinEvent.Message = strings.Replace(trimmedMessage, nbsp, " ", -1)

	user, err := model.GetOrCreateUserBySlackID(slackId)

	if err != nil {
		log.Printf("[ERROR] Could not find or create a user for slack id %s", slackId)
		return &coinEvent, err
	}

	coinEvent.User = user

	log.Printf("Created CoinEvent: message - %s", coinEvent.Message)
	return &coinEvent, nil
}
