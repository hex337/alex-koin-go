package main

// Here lives database migrations
// Heroku runs this automagically before the server command (see Procfile)

import (
	"log"

	"github.com/hex337/alex-koin-go/Config"
	"github.com/hex337/alex-koin-go/Model"
)

func main() {

	log.Println("Starting database migrations")
	Config.GetEnvVars()
	Config.DBOpen()

	// TODO SELECT pg_try_advisory_lock(migration); and if f don't run migration as they are already happening
	Config.DB.AutoMigrate(&Model.Transaction{}, &Model.User{}, &Model.Coin{})
	// TODO SELECT pg_advisory_unlock(migration);

}
