package main

import (
	"github.com/hex337/alex-koin-go/Config"
	"github.com/hex337/alex-koin-go/Model"

	"log"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	var err error
	Config.GetEnvVars()
	Config.DBOpen()
	// defer Config.DB.Close()

	Config.DB.Migrator().DropTable(&Model.User{})
	Config.DB.Migrator().DropTable(&Model.Coin{})
	Config.DB.Migrator().DropTable(&Model.Transaction{})

	Config.DB.AutoMigrate(&Model.Transaction{}, &Model.User{}, &Model.Coin{})

	user1 := &Model.User{ FirstName: "Alex", LastName: "Koin" }
	user2 := &Model.User{ FirstName: "Koin", LastName: "Lord" }

	err =  Model.CreateUser(user1)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	err =  Model.CreateUser(user2)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	coin := &Model.Coin{ 
		Origin: "Reason for the season",
		MinedByUserID: user1.Model.ID,
		UserID: user2.Model.ID,
		CreatedByUserId: user1.Model.ID,
		Transactions: []Model.Transaction{
			{
				Amount:     1,
				Memo:       "Initial Koin Creation",
				FromUserID: user1.Model.ID,
				ToUserID:   user2.Model.ID,
			},
		},
	}

	err =  Model.CreateCoin(coin)
	if err != nil {
		log.Fatalf("Could not create coin : %s", err.Error())
	}

	var user Model.User
	err = Model.GetUserByID(&user, "1")
	if err != nil {
		log.Println(err.Error())
	}

	spew.Dump(user1)
	log.Printf("\nUser 1 Balance : %d\nUser 2 Balance : %d", 
		Model.GetUserBalance(user1),
		Model.GetUserBalance(user2),
	)
}
