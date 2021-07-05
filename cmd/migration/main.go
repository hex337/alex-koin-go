package main

// Here lives database migrations
// Heroku runs this automagically before the server command (see Procfile)

import (
	"log"

	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"
)

func main() {

	log.Println("Starting database migrations")
	config.GetEnvVars()
	config.DBOpen()

	// TODO SELECT pg_try_advisory_lock(migration); and if f don't run migration as they are already happening
	config.DB.AutoMigrate(&model.User{}, &model.Coin{}, &model.Transaction{})
	// TODO SELECT pg_advisory_unlock(migration);

}
