package main

import (
	"log"
	"net/http"

	"github.com/hex337/alex-koin-go/Config"
	"github.com/hex337/alex-koin-go/Router"
)

func main() {
	Config.GetEnvVars()
	Config.DBOpen()

	Router.SlackEvents()

	serverURL := Config.ServerURL(Config.BuildServerConfig())
	log.Printf("[INFO] Server listening %s", serverURL)
	http.ListenAndServe(serverURL, nil)
}
