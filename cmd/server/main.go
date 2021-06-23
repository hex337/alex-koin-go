package main

import (
	"log"
	"net/http"

	"github.com/hex337/alex-koin-go/Config"
	"github.com/hex337/alex-koin-go/Models"
	"github.com/hex337/alex-koin-go/endpoints"
)

func main() {
	var err error

	Config.GetEnvVars()
	Config.DBOpen()

	// defer config.DB.Close()
	Config.DB.AutoMigrate(&Models.User{}, &Models.Coin{}, &Models.Transaction{})

	var user []Models.User
	err = Models.GetAllUsers(&user)
	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("%v", user[0])

	endpoints.SlackEvents()

	serverURL := Config.ServerURL(Config.BuildServerConfig())
	log.Printf("[INFO] Server listening %s", serverURL)
	http.ListenAndServe(serverURL, nil)
}
