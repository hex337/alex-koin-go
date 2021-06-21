package main

import (
	"log"
	"net/http"
	"os"
	// "time"

	"github.com/hex337/alex-koin-go/endpoints"
	"github.com/hex337/alex-koin-go/Models"
	"github.com/hex337/alex-koin-go/Config"
	"github.com/joho/godotenv"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

func main() {
	var err error
	_, skipEnvFile := os.LookupEnv("SKIP_ENV_FILE")
	if !skipEnvFile {
		err = godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %s", err.Error())
		}
	}

	Config.DB, err = gorm.Open(postgres.Open(Config.DBURL(Config.BuildDBConfig())))
  if err != nil {
		log.Fatalf("Could not connect to db : %s", err.Error())
	}

	// defer config.DB.Close()
	// config.DB.AutoMigrate(&models.User{})

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
