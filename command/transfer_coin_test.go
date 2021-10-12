package command

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseMessageWithSingleReceiver(t *testing.T) {
	numOfCoins := "10"
	receiverId := "U000AAA"
	reason := "being awesome"

	parsedResult := parseMessage(fmt.Sprintf("transfer %s to <@%s> %s", numOfCoins, receiverId, reason))

	assert.Equal(t, numOfCoins, parsedResult["amount"], "matches coin")
	assert.Equal(t, receiverId, parsedResult["to_slack_ids"], "matches receiverId")
	assert.Equal(t, reason, parsedResult["reason"], "matches reason")
}

func TestParseMessageWithMultipleReceivers(t *testing.T) {
	numOfCoins := "10"
	receiverId := "U000AAA"
	secondReceiverId := "U000BBB"
	reason := "being awesome"

	parsedResult := parseMessage(fmt.Sprintf("transfer %s to <@%s> <@%s> %s", numOfCoins, receiverId, secondReceiverId, reason))

	assert.Equal(t, numOfCoins, parsedResult["amount"], "matches coin")
	assert.Equal(t, fmt.Sprintf("%s,%s", receiverId, secondReceiverId), parsedResult["to_slack_ids"], "matches receiverId")
	assert.Equal(t, reason, parsedResult["reason"], "matches reason")
}

func TestSplitCoinsWithSingleReceiver(t *testing.T) {
	coins := splitCoins(5, 3)

	sum := 0
	for _, coin := range coins {
		assert.LessOrEqual(t, coin, 5, "each coin amount is <= max")
		sum += coin
	}

	assert.Equal(t, 5, sum, "amount is the same")
}

func TestSplitCoinsWithMultipleReceivers(t *testing.T) {
	coins := splitCoins(5, 1)

	sum := 0
	for _, coin := range coins {
		assert.LessOrEqual(t, 5, coin, "each coin amount is <= max")
		sum += coin
	}

	assert.Equal(t, 5, sum, "amount is the same")
}