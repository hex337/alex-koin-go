package Config

import (
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func DBOpen() {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:              time.Second,   // Slow SQL threshold
			LogLevel:                   logger.Info, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,          // Disable color
		},
	)

	DB, err = gorm.Open(
		postgres.Open(dbUrl()),
		&gorm.Config{ Logger: newLogger },
	)
  if err != nil {
		log.Fatalf("Could not connect to db : %s", err.Error())
	}
}

func dbUrl() string {
	env, present := os.LookupEnv("DATABASE_URL")

	if !present {
		log.Fatal("DATABASE_URL must be set as an environment variable")
	}
	return env
}
