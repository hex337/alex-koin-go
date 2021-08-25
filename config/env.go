package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func GetNoPrefixChannelIds() map[string]int {
	var idMap map[string]int
	idMap = make(map[string]int)
	channelIds, _ := os.LookupEnv("CHANNEL_IDS")
	// if !channelIdsProvided {
	// 	log.Fatalf("env var CHANNEL_IDS not set")
	// }

	for _, element := range strings.Split(channelIds, ",") {
		idMap[element] = 1
	}

	return idMap
}

func GetCurrencyShortName() string {
	currencyShortName, envVarProvided := os.LookupEnv("CURRENCY_SHORT_NAME")
	if !envVarProvided {
		log.Fatalf("env var CURRENCY_SHORT_NAME not set")
	}
	return currencyShortName
}

func GetCurrencyName() string {
	currencyName, envVarProvided := os.LookupEnv("CURRENCY_NAME")
	if !envVarProvided {
		log.Fatalf("env var CURRENCY_NAME not set")
	}
	return currencyName
}

func GetBotSlackID() string {
	botId, envVarProvided := os.LookupEnv("SLACK_BOT_ID")
	if !envVarProvided {
		log.Fatalf("env var SLACK_BOT_ID not set")
	}
	return botId
}

func GetAdminSlackIds() map[string]int {
	var idMap map[string]int
	idMap = make(map[string]int)
	adminIds, adminIdsProvided := os.LookupEnv("ADMIN_IDS")
	if !adminIdsProvided {
		log.Fatalf("env var ADMIN_IDS not set")
	}

	for _, element := range strings.Split(adminIds, ",") {
		idMap[element] = 1
	}

	return idMap
}

func GetKoinLordSlackIds() map[string]int {
	var idMap map[string]int
	idMap = make(map[string]int)
	koinLordIds, koinLordIdsProvided := os.LookupEnv("KOIN_LORD_IDS")
	if !koinLordIdsProvided {
		log.Fatalf("env var KOIN_LORD_IDS not set")
	}

	for _, element := range strings.Split(koinLordIds, ",") {
		idMap[element] = 1
	}

	return idMap
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
