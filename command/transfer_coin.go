package command

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hex337/alex-koin-go/config"
	"github.com/hex337/alex-koin-go/model"

	"log"
	"regexp"
)

type TransferCoinCommand struct{}

func (c *TransferCoinCommand) Run(msg string, event *CoinEvent) (string, error) {
	parsedResult, err := parseMessage(event.Message)
	if err != nil {
		return err.Error(), nil
	}

	// parsing error should be covered in func validateMessage
	totalAmount, _ := strconv.Atoi(parsedResult["amount"])
	toUserIds := strings.Split(parsedResult["to_slack_ids"], ",")
	reason := parsedResult["reason"]
	splitAmounts := splitCoins(totalAmount, len(toUserIds))

	for i, toUserId := range toUserIds {
		amount := splitAmounts[i]
		log.Printf("INF transfer coin amt: %d, toUserIds: %s, reason: %s", totalAmount, toUserIds, reason)

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

		log.Printf("Transfered %d koin from %d to %s.", amount, event.User.ID, toUserId)
	}

	return fmt.Sprintf("Transfered %d koin.", totalAmount), nil
}

func parseMessage(message string) (map[string]string, error) {
	re := regexp.MustCompile(`^(?i)transfer (?P<amount>[0-9]+) (?:to )?(?P<to_slack_ids>\<[@0-9A-Z<> ]+\>) (?:for)?(?P<reason>.+)`)
	matches := re.FindStringSubmatch(message)

	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}

	slackIdsRe := regexp.MustCompile(`(<@(?P<to_slack_id>[0-9A-Z]+)>)+`)
	slackIdMatches := slackIdsRe.FindAllStringSubmatch(result["to_slack_ids"], -1)
	slackIds := make([]string, len(slackIdMatches))

	for i, name := range slackIdsRe.SubexpNames() {
		if name != "to_slack_id" {
			continue
		}

		for j, v := range slackIdMatches {
			slackIds[j] = v[i]
		}
	}
	result["to_slack_ids"] = strings.Join(slackIds, ",")

	err := validateMessage(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func validateMessage(parsedResult map[string]string) error {
	requiredKeys := []string {
		"amount",
		"to_slack_ids",
		"reason",
	}

	for _, key := range requiredKeys {
		value, ok := parsedResult[key]
		if !ok || len(value) == 0 {
			return errors.New("Invalid transfer format. Expected something like `@Alex Koin transfer 3 to @alexk for being amazing`. See the channel details for command syntax.")
		}
	}

	amount, err := strconv.Atoi(parsedResult["amount"])
	if err != nil {
		log.Printf("INF amount not parsed as int: %s", parsedResult["amount"])
		return errors.New("Invalid transfer amount.")
	}

	numOfReceivers := len(strings.Split(parsedResult["to_slack_ids"], ","))
	if amount < numOfReceivers {
		log.Printf("amount is not enough for split: amount: %d, receivers: %d", amount, numOfReceivers)
		return errors.New(fmt.Sprintf("alex koin does not support fraction"))
	}

	return nil
}

func splitCoins(amount int, count int) []int {
	rand.Seed(time.Now().UnixNano())
	coins := make([]int, count)
	leftCount := count - 1
	for i, _ := range coins {
		coins[i] = rand.Intn(amount - leftCount - 1) + 1
		leftCount -= 1
		amount -= coins[i]
	}
	coins[rand.Intn(len(coins))] += amount
	return coins
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
