package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hex337/alex-koin-go/endpoints"

	"github.com/joho/godotenv"
)

// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"

func main() {
	_, skipEnvFile := os.LookupEnv("SKIP_ENV_FILE")
	if !skipEnvFile {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	endpoints.SlackEvents()

	log.Println("[INFO] Server listening")

	port, portProvided := os.LookupEnv("PORT")
	if !portProvided {
		port = "3000"
	}
	http.ListenAndServe(":"+port, nil)
}
