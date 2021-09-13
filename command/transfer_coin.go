package command

import (
	"fmt"
	"strconv"

	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"

	"log"
	"regexp"
)

type TransferCoinCommand struct{}

func (c *TransferCoinCommand) Run(msg string, event *CoinEvent) (string, error) {
	re := regexp.MustCompile(`^(?i)transfer (?P<amount>[0-9]+) (?:to )?\<@(?P<to_slack_id>[0-9A-Z]+)\> (?:for)??(?P<reason>.+)`)
	matches := re.FindStringSubmatch(event.Message)

	if len(matches) < 4 {
		return "Invalid transfer format. Expected something like `@Alex Koin transfer 3 to @alexk for being amazing`. See the channel details for command syntax.", nil
	}

	amount, err := strconv.Atoi(matches[1])
	if err != nil {
		log.Printf("INF amount not parsed as int: %s", matches[1])
		return "Invalid transfer amount.", nil
	}
	toUserId := matches[2]
	reason := matches[3]

	log.Printf("INF transfer coin amt: %d, toUserId: %s, reason: %s", amount, toUserId, reason)

	toUser, err := model.GetOrCreateUserBySlackID(toUserId)

	if err != nil {
		log.Printf("Could not find user with slack id %s: %s", toUserId, err.Error())
		return "", err
	}

	canTransfer, msg := canTransferCoin(event.User, toUser, amount)

	if !canTransfer {
		return msg, nil
	}

	var coinIds []int
	err = config.DB.Table("coins").Select("id").Where("user_id = ?", event.User.ID).Limit(amount).Find(&coinIds).Error
	if err != nil {
		log.Printf("ERR error fetching coin ids: %s", err)
		return "So uh... something went wrong. D'oh.", err
	}

	err = config.DB.Table("coins").Where("user_id = ? AND id IN ?", event.User.ID, coinIds).UpdateColumn("user_id", toUser.ID).Error
	if err != nil {
		log.Printf("ERR error updating coins: %s", err)
		return "So uh... something went wrong. D'oh.", err
	}

	transfer := &model.Transaction{
		Amount:     amount,
		Memo:       reason,
		FromUserID: event.User.ID,
		ToUserID:   toUser.ID,
	}

	err = model.CreateTransaction(transfer)

	if err != nil {
		log.Printf("Could not transfer coin(s) : %s", err.Error())
		return "", err
	}

	return fmt.Sprintf("Transfered %d koin.", amount), nil
}

func canTransferCoin(sender *model.User, receiver *model.User, amount int) (bool, string) {
	if sender.ID == receiver.ID {
		return false, "This action is very :sus:. We have notified the Lords of Koin about your behavior."
	}

	if amount <= 0 {
		return false, "How about a big ol' ball of nope."
	}

	if sender.GetBalance() < int64(amount) {
		return false, "You lack the koin for such a transfer."
	}

	return true, ""
}
