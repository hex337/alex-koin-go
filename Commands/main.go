package Commands

import (
	"strings"
	"os"

	"github.com/slack-go/slack"
)

func ProcessMessage(channel string, msgTimestamp string, msg string) (err error) {

	lowerCaseMsg := strings.ToLower(msg)

	if (strings.HasPrefix(lowerCaseMsg, "my balance")) {
		replyWith(channel, msgTimestamp, "Your Broke")
	}

  return err
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
