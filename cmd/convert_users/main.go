package main

import (
	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"

	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/slack-go/slack"
)

func main() {

	config.GetEnvVars()
	config.DBOpen()

	var users []model.User
	res := config.DB.Where("slack_id ~ ?", "^U.*").Find(&users)

	log.Printf("Found %d", res.RowsAffected)

	for _, user := range users {
		migrate(user.SlackID)
	}
}

func migrate(id string) {

	botSecret := os.Getenv("SLACK_BOT_SECRET")
	var api = slack.New(botSecret)

	user, err := api.GetUserInfo(id)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(user)
	log.Fatal("Start with one user for now")
}
