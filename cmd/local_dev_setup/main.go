package main

import (
	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"

	"log"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	var err error
	config.GetEnvVars()
	config.DBOpen()
	// defer Config.DB.Close()

	config.DB.Migrator().DropTable(&model.User{})
	config.DB.Migrator().DropTable(&model.Coin{})
	config.DB.Migrator().DropTable(&model.Transaction{})

	config.DB.AutoMigrate(&model.Transaction{}, &model.User{}, &model.Coin{})

	user1 := &model.User{FirstName: "Alex", LastName: "Koin"}
	user2 := &model.User{FirstName: "Koin", LastName: "Lord", SlackID: "W0122R46YBC"}

	err = model.CreateUser(user1)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	err = model.CreateUser(user2)
	if err != nil {
		log.Fatalf("Could not create user : %s", err.Error())
	}

	coin := &model.Coin{
		Origin:          "Reason for the season",
		MinedByUserID:   user1.Model.ID,
		UserID:          user2.Model.ID,
		CreatedByUserId: user1.Model.ID,
		Transactions: []model.Transaction{
			{
				Amount:     1,
				Memo:       "Initial Koin Creation",
				FromUserID: user1.Model.ID,
				ToUserID:   user2.Model.ID,
			},
		},
	}

	err = model.CreateCoin(coin)
	if err != nil {
		log.Fatalf("Could not create coin : %s", err.Error())
	}

	user1, res, err := model.GetUserByID(1)
	if err != nil && res == false {
		log.Println(err.Error())
	}

	spew.Dump(user1)
	log.Printf("\nUser 1 Balance : %d\nUser 2 Balance : %d",
		user1.GetBalance(),
		user2.GetBalance(),
	)
}
