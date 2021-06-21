
package main

import (
	"log"
	"os"
	"time"

	"github.com/hex337/alex-koin-go/Models"
	"github.com/hex337/alex-koin-go/Config"
	"github.com/joho/godotenv"
	"github.com/davecgh/go-spew/spew"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:              time.Second,   // Slow SQL threshold
			LogLevel:                   logger.Info, // Log level
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,          // Disable color
		},
	)

	Config.DB, err = gorm.Open(
		postgres.Open(Config.DBURL(Config.BuildDBConfig())),
		&gorm.Config{ Logger: newLogger },
	)
  if err != nil {
		log.Fatalf("Could not connect to db : %s", err.Error())
	}

  

	// defer Config.DB.Close()

	Config.DB.Migrator().DropTable(&Models.User{})
	Config.DB.Migrator().DropTable(&Models.Coin{})
	Config.DB.Migrator().DropTable(&Models.Transaction{})

	Config.DB.AutoMigrate(&Models.Transaction{}, &Models.User{}, &Models.Coin{})

	user1 := &Models.User{ FirstName: "Alex", LastName: "Koin" }
	user2 := &Models.User{ FirstName: "Koin", LastName: "Lord" }

	err =  Models.CreateUser(user1)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	err =  Models.CreateUser(user2)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	coin := &Models.Coin{ 
		Origin: "Reason for the season",
		MinedByUserID: user1.Model.ID,
		UserID: user2.Model.ID,
		CreatedByUserId: user1.Model.ID,
		Transactions: []Models.Transaction{
			{
				Amount: 1,
				Memo: "Initial Koin Creation",
				FromUserID: user1.Model.ID,
				ToUserID: user2.Model.ID,
		  },
	  },
  }

	err =  Models.CreateCoin(coin)
	if err != nil {
		log.Fatalf("Could not create coin : %s", err.Error())
	}

	var user Models.User
	err = Models.GetUserByID(&user, "1")
	if err != nil {
		log.Println(err.Error())
	}

	spew.Dump(user1)
	log.Printf("\nUser 1 Balance : %d\nUser 2 Balance : %d", 
		Models.GetUserBalance(user1),
		Models.GetUserBalance(user2),
	)
}
