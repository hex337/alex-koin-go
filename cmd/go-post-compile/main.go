package main

// Here lives database migrations
// I hate the name too. This is called by the heroku buildpack before the procfile is booted

import (
	"github.com/hex337/alex-koin-go/Config"
	"github.com/hex337/alex-koin-go/Model"
)

func main() {

	Config.GetEnvVars()
	Config.DBOpen()

	Config.DB.AutoMigrate(&Model.Transaction{}, &Model.User{}, &Model.Coin{})
}
