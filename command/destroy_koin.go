package command

import (
	"fmt"

	"github.com/hex337/alex-koin-go/model"

	"log"
	"regexp"
)

type DestroyKoinCommand struct{}

func (c *DestroyKoinCommand) Run(msg string, event *CoinEvent) (string, error) {
	re := regexp.MustCompile(`^(?i)destroy koin (?:for )?\<@(?P<to_slack_id>[0-9A-Z]+)\> (?:for)??(?P<reason>.+)`)
	matches := re.FindStringSubmatch(event.Message)

	if len(matches) < 3 {
		return "Invalid destroy koin format. Expected something like `@Alex Koin destroy koin @alexk for being a bad cat`. See the channel details for command syntax.", nil
	}

	toUserId := matches[1]
	reason := matches[2]

	// Get or create the to user
	toUser, err := model.GetOrCreateUserBySlackID(toUserId)

	if err != nil {
		log.Fatalf("Could not find user with slack id %s: %s", toUserId, err.Error())
		return "", err
	}

	canDestroy, msg := canDestroyCoin(event.User, toUser)

	if !canDestroy {
		return msg, nil
	}

	coin, err := toUser.GetCoin()

	if err != nil {
		log.Fatalf("Failed to Destroy Koin. Err: %s", err)
		return fmt.Sprintf("There was a failure in the system. Koin not destroyed"), nil
	}

	err = coin.DestroyKoin()

	if err != nil {
		log.Fatalf("Could not properly destroy koin. Err: %s", err)
		return "Koin failed to be destroyed.", nil
	}

	return fmt.Sprintf("How terribly unfortunate. A koin has been destroyed because: %s. Do honor your Lords.", reason), nil
}

/**
 * Rules:
 *  Only Admin and Lords can destroy Koin
 */
 func canDestroyCoin(sender *model.User, receiver *model.User) (bool, string) {
	role := sender.Role()

	if receiver.GetBalance() == 0 {
		return false, "Have pitty on this poor soul, for they've no coin left to destroy"
	}

	if role.Admin || role.Lord {
		return true, ""
	}

	return false, "You pathetic human you don't have any family any friends or any land."
 }
