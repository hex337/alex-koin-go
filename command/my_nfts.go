package command

import (
	"fmt"
	"log"

	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"
)

type MyNftsCommand struct{}

func (c *MyNftsCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
	var nfts []model.Nft

	nft_result := config.DB.Table("nfts").Where("owned_by_user_id = ?", event.User.ID).Find(&nfts)
	if nft_result.Error != nil {
		msg = fmt.Sprintf("Issue fetching NFTs for user: %s", nft_result.Error)
		log.Printf("ERR: %s", msg)
		return BotResponse{Text: msg}, nil
	}

	var response_msg string

	if len(nfts) > 0 {
		response_msg = "Your NFTs:```"

		for _, nft := range nfts {
			response_msg += fmt.Sprintf("Name: %15s | Paid: %3d | Link: %s\n", nft.Name, nft.PricePaid, nft.DisplayURL)
		}

		response_msg += "```"
	} else {
		response_msg = "You don't currently own any NFTs."
	}

	return BotResponse{Text: response_msg}, nil
}
