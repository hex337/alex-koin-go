package main

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/hex337/alex-koin-go/Config"
	"github.com/hex337/alex-koin-go/Models"
)

func main() {
	var err error
	Config.GetEnvVars()
	Config.DBOpen()
	// defer Config.DB.Close()

	Config.DB.Migrator().DropTable(&Models.User{})
	Config.DB.Migrator().DropTable(&Models.Coin{})
	Config.DB.Migrator().DropTable(&Models.Transaction{})

	Config.DB.AutoMigrate(&Models.User{}, &Models.Coin{}, &Models.Transaction{})

	user1 := &Models.User{FirstName: "Alex", LastName: "Koin"}
	user2 := &Models.User{FirstName: "Koin", LastName: "Lord"}

	err = Models.CreateUser(user1)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	err = Models.CreateUser(user2)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	coin := &Models.Coin{
		Origin:          "Reason for the season",
		MinedByUserID:   user1.Model.ID,
		UserID:          user2.Model.ID,
		CreatedByUserId: user1.Model.ID,
		Transactions: []Models.Transaction{
			{
				Amount:     1,
				Memo:       "Initial Koin Creation",
				FromUserID: user1.Model.ID,
				ToUserID:   user2.Model.ID,
			},
		},
	}

	err = Models.CreateCoin(coin)
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
