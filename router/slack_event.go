package router

import (
	"log"

	"github.com/hex337/alex-koin-go/command"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func SlackEvents() {
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")

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
			switch ev := innerEvent.Data.(type) {

			// case *slackevents.AppMentionEvent:
			// 	log.Printf("INF Recvd AppMentionedEvent %+v", ev)
			// 	go command.ProcessAppMentionEvent(ev)
			case *slackevents.MessageEvent:
				log.Printf("INF Recvd AppMentionedEvent %+v", ev)
				go command.ProcessMessageEvent(ev)
			}
		}
	})
}
