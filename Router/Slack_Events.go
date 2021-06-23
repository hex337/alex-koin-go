package Router

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func SlackEvents () {
	botSecret := os.Getenv("SLACK_BOT_SECRET")
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")

	var api = slack.New(botSecret)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			log.Printf("%v\n", innerEvent.Data)
			switch ev := innerEvent.Data.(type) {

			// Someone mentioning the bot by name
			case *slackevents.AppMentionEvent:
				log.Printf("%#v", ev)
				text := strings.Split(ev.Text, " ")[0]
				api.PostMessage(
					ev.Channel,
					slack.MsgOptionText(text, false),
					slack.MsgOptionTS(ev.TimeStamp), // reply in thread
				)

			}
		}
	})
}
