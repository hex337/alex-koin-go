package command

import (
	"fmt"

	"github.com/hex337/alex-koin-go/model"

	"log"
	"regexp"
)

type CreateNftCommand struct{}

func (c *CreateNftCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
	re := regexp.MustCompile(`^(?i)create nft (?P<name>[0-9a-zA-Z'-_]+) (?P<url>[><0-9a-zA-Z.:\/?&%]+)`)
	matches := re.FindStringSubmatch(event.Message)

	if len(matches) < 2 {
		return BotResponse{Text: "Invalid create nft format. Expected something like `@Alex Koin create nft nft name http://url.com`. See the channel details for command syntax."}, nil
	}

	nft_name := matches[1]
	source_url := matches[2]

	canCreate, msg := canCreateNft(event.User)
	if !canCreate {
		return BotResponse{Text: msg}, nil
	}

	nft := &model.Nft{
		OwnedByUserId:   0, // means this is not owned by anyone
		CreatedByUserId: event.User.ID,
		MinBid:          0,
		PricePaid:       0, // if price paid is 0, that means it was just created
		Name:            nft_name,
		SourceURL:       source_url,
		DisplayURL:      source_url,
	}

	err := model.CreateNft(nft)
	if err != nil {
		log.Printf("Could not create nft : %s", err.Error())
		return BotResponse{Text: ""}, err
	}

	return BotResponse{Text: fmt.Sprintf("Created new nft `%s:%s` with source `%s`.", nft.Name, nft.Hash, nft.SourceURL)}, nil
}

func canCreateNft(creator *model.User) (bool, string) {
	role := creator.Role()

	if role.Admin {
		return true, ""
	}

	return false, "Sorry buddy, only Admins can create nfts right now."
}
