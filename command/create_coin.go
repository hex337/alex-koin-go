package command

import (
	"fmt"

	"github.com/hex337/alex-koin-go/model"

	"log"
	"regexp"
)

type CreateCoinCommand struct{}

func (c *CreateCoinCommand) Run(msg string, event *CoinEvent) (BotResponse, error) {
	re := regexp.MustCompile(`^(?i)create koin (?:for )?\<@(?P<to_slack_id>[0-9A-Z]+)\> (?:for)??(?P<reason>.+)`)
	matches := re.FindStringSubmatch(event.Message)

	if len(matches) < 3 {
		return BotResponse{Text: "Invalid create koin format. Expected something like `@Alex Koin create koin @alexk for being amazing`. See the channel details for command syntax."}, nil
	}

	toUserId := matches[1]
	reason := matches[2]

	// Get or create the to user
	toUser, err := model.GetOrCreateUserBySlackID(toUserId)

	if err != nil {
		log.Printf("Could not find user with slack id %s: %s", toUserId, err.Error())
		return BotResponse{Text: ""}, err
	}

	canCreate, msg := canCreateCoin(event.User, toUser)

	if !canCreate {
		return BotResponse{Text: msg}, nil
	}

	coin := &model.Coin{
		Origin:          reason,
		MinedByUserID:   event.User.ID,
		UserID:          toUser.ID,
		CreatedByUserId: event.User.ID,
	}

	err = model.CreateCoin(coin)
	if err != nil {
		log.Printf("Could not create coin : %s", err.Error())
		return BotResponse{Text: ""}, err
	}

	transaction := &model.Transaction{
		Amount:     1,
		Memo:       "Initial Coin Creation",
		FromUserID: event.User.ID,
		ToUserID:   toUser.ID,
	}

	err = model.CreateTransaction(transaction)
	if err != nil {
		log.Printf("Could not create transaction : %s", err.Error())
		return BotResponse{Text: ""}, err
	}

	return BotResponse{Text: fmt.Sprintf("Created new koin `%s` with reason '%s'.", coin.Hash, reason)}, nil
}

/**
 * Rules:
 *  Admin can create coins no matter what
 *  Lord can create as many coins as wanted, except for themselves
 *  Anyone else is one coin per week no matter what
 */
func canCreateCoin(sender *model.User, receiver *model.User) (bool, string) {
	role := sender.Role()

	if role.Admin {
		return true, ""
	}

	// Only admins can create for themselves
	if sender.ID == receiver.ID {
		return false, "How about a big ol' ball of nope on that one."
	}

	if role.Lord {
		return true, ""
	}

	coinCount := model.CoinsCreatedThisWeekForUser(sender)

	if coinCount < 1 {
		return true, ""
	}

	return false, "You can only create one koin per week, try again on Monday!"
}
