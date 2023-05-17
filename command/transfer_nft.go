package command

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"
)

type TransferNftCommand struct{}

func (c *TransferNftCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
	re := regexp.MustCompile(`^(?i)transfer nft (?P<nft_name>[0-9a-zA-z'-_]+) (?:to )?\<@(?P<to_slack_id>[0-9A-Z]+)\> (?:for )??(?P<amount>[0-9]+)`)
	matches := re.FindStringSubmatch(event.Message)

	if len(matches) < 4 {
		return BotResponse{Text: "Invalid transfer format. Expected something like `@Alex Koin transfer nft cool_gif to @alexk for 5`. See the channel details for command syntax."}, nil
	}

	log.Printf("INFO: Matches %s", matches)
	log.Printf("INFO: matches[0]: %s", matches[0])
	log.Printf("INFO: matches[1]: %s", matches[1])
	log.Printf("INFO: matches[2]: %s", matches[2])
	log.Printf("INFO: matches[3]: %s", matches[3])
	amount, err := strconv.Atoi(matches[3])
	amt := uint(amount)
	if err != nil {
		log.Printf("INF amount not parsed as int: %s", matches[3])
		return BotResponse{Text: "Invalid transfer amount."}, nil
	}
	toUserId := matches[2]
	nft_name := matches[1]

	log.Printf("INF transfer nft %s for %d, toUserId: %s", nft_name, amt, toUserId)

	toUser, err := model.GetOrCreateUserBySlackID(toUserId)
	if err != nil {
		log.Printf("Could not find user with slack id %s: %s", toUserId, err.Error())
		return BotResponse{Text: ""}, err
	}

	var nft model.Nft
	nft_result := config.DB.Table("nfts").Where("name = ?", nft_name).Find(&nft)
	if nft_result.Error != nil {
		msg = fmt.Sprintf("Could not find nft with name %s", nft_name)
		log.Printf("INFO: %s", msg)
		return BotResponse{Text: msg}, nil
	}

	role := event.User.Role()
	if !role.Admin && nft.OwnedByUserId != event.User.ID {
		msg = fmt.Sprintf("You don't own this nft, so you can't trade it.")
		log.Printf("INFO: %s", msg)
		return BotResponse{Text: msg}, nil
	}

	canTransfer, msg := canTransferNft(event.User, toUser, amount)
	if !canTransfer {
		return BotResponse{Text: msg}, nil
	}

	// transfer the coins and then transfer the nft
	var coinIds []int
	err = config.DB.Table("coins").Select("id").Where("user_id = ?", toUser.ID).Limit(amount).Find(&coinIds).Error
	if err != nil {
		log.Printf("ERR error fetching coin ids: %s", err)
		return BotResponse{Text: "So uh... something went wrong. D'oh."}, err
	}

	err = config.DB.Table("coins").Where("user_id = ? AND id IN ?", toUser.ID, coinIds).UpdateColumn("user_id", event.User.ID).Error
	if err != nil {
		log.Printf("ERR error updating coins: %s", err)
		return BotResponse{Text: "So uh... something went wrong. D'oh."}, err
	}

	config.DB.Model(&nft).Select("OwnedByUserId", "PricePaid").Updates(model.Nft{OwnedByUserId: toUser.ID, PricePaid: amt})
	if err != nil {
		log.Printf("ERR error moving the NFT: %s", err)
		return BotResponse{Text: "So uh... something went wrong. D'oh."}, err
	}

	return BotResponse{Text: "Transferred the NFT."}, nil
}

func canTransferNft(sender *model.User, receiver *model.User, amount int) (bool, string) {
	role := sender.Role()

	if receiver.GetBalance() < int64(amount) {
		return false, "Not enough koin for this transaction."
	}

	if role.Admin {
		return true, ""
	}

	if sender.ID == receiver.ID {
		return false, "This action is very :sus:. We have notified the Lords of Koin about your behavior."
	}

	if amount <= 0 {
		return false, "This is pretty much a great example of a nope rope. Nope nope nope."
	}

	return true, ""
}
