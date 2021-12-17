package command

import (
	"fmt"
	"log"

	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"
)

type AllNftsCommand struct{}

func (c *AllNftsCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
	var nfts []model.Nft

	nft_result := config.DB.Table("nfts").Where("owned_by_user_id > 0").Find(&nfts)
	if nft_result.Error != nil {
		msg = fmt.Sprintf("Issue fetching all NFTs: %s", nft_result.Error)
		log.Printf("ERR: %s", msg)
		return BotResponse{Text: msg}, nil
	}

	var response_msg string

	if len(nfts) > 0 {
		response_msg = "All NFTs:```"

		for _, nft := range nfts {
			var owner model.User
			err := config.DB.Table("users").Where("ID = ?", nft.OwnedByUserId).Limit(1).Find(&owner).Error
			if err != nil {
				log.Printf("ERR: error getting the owner of an nft: %s", err)
				return BotResponse{Text: "Error fetching all NFTs..."}, nil

			}
			response_msg += fmt.Sprintf("Name: %15s | Owner: %15s | Paid: %3d | Link: %s\n", nft.Name, owner.FirstName+" "+owner.LastName, nft.PricePaid, nft.DisplayURL)
		}

		response_msg += "```"
	} else {
		response_msg = "You don't currently own any NFTs."
	}

	return BotResponse{Text: response_msg}, nil
}
