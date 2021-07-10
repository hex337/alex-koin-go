package command

import (
	"fmt"

	"github.com/hex337/alex-koin-go/model"

	"log"
	"regexp"
)

type CreateCoinCommand struct{}

func (c *CreateCoinCommand) Run(msg string, event *CoinEvent) (string, error) {
	re := regexp.MustCompile(`^(?i)create koin \<@(?P<to_slack_id>[0-9A-Z]+)\> (?:for)??(?P<reason>.+)`)
	matches := re.FindStringSubmatch(event.Message)

	toUserId := matches[1]
	reason := matches[2]

	// Get or create the to user
	toUser, err := model.GetOrCreateUserBySlackID(toUserId)

	canCreate, msg := canCreateCoin(event.User, toUser)

	if !canCreate {
		return msg, nil
	}

	if err != nil {
		log.Fatalf("Could not find user with slack id %s: %s", toUserId, err.Error())
		return "", err
	}

	coin := &model.Coin{
		Origin:          reason,
		MinedByUserID:   event.User.ID,
		UserID:          toUser.ID,
		CreatedByUserId: event.User.ID,
		Transactions: []model.Transaction{
			{
				Amount:     1,
				Memo:       "Initial Coin Creation",
				FromUserID: event.User.ID,
				ToUserID:   toUser.ID,
			},
		},
	}

	err = model.CreateCoin(coin)

	if err != nil {
		log.Fatalf("Could not create coin : %s", err.Error())
		return "", err
	}

	return fmt.Sprintf("Created new koin `%s` with reason '%s'.", coin.Hash, reason), nil
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

	return false, "Only supported for Admins and Lords."
}
