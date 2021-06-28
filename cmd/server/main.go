package main

import (
	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/router"

	"log"
	"net/http"
)

func main() {
	config.GetEnvVars()
	config.DBOpen()

	router.SlackEvents()

	serverURL := config.ServerURL(config.BuildServerConfig())
	log.Printf("[INFO] Server listening %s", serverURL)
	http.ListenAndServe(serverURL, nil)
}
