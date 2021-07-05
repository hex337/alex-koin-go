package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetBotSlackID() string {
	botId, botIdProvided := os.LookupEnv("SLACK_BOT_ID")
	if !botIdProvided {
		log.Fatalf("env var SLACK_BOT_ID not set")
	}
	return botId
}

func GetEnvVars() {
	var err error
	_, skipEnvFile := os.LookupEnv("SKIP_ENV_FILE")
	if !skipEnvFile {
		err = godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err.Error())
		}
	}
}
