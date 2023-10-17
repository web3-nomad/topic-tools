package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
	"strings"
)

func main() {
	var client *hedera.Client
	var err error

	client, err = hedera.ClientForName(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		println(err.Error(), ": error creating client")
		return
	}
	operatorAccountID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to AccountID")
		return
	}

	// Retrieving operator key from environment variable OPERATOR_KEY
	//	operatorKey, err := hedera.PrivateKeyFromSeedEd25519(os.Getenv("OPERATOR_KEY"))

	mnemonic, err := hedera.NewMnemonic(strings.Split(os.Getenv("OPERATOR_SEED"), " "))
	if err != nil {
		println(err.Error(), ": error converting string to mnemonic")
		return
	}
	operatorKey, err := mnemonic.ToStandardEd25519PrivateKey("", 0)
	if err != nil {
		println(err.Error(), ": error converting string to PrivateKey")
		return
	}

	// Defaults the operator account ID and key such that all generated transactions will be paid for
	// by this account and be signed by this key
	client.SetOperator(operatorAccountID, operatorKey)

	// Make a new topic
	transactionResponse, err := hedera.NewTopicCreateTransaction().
		SetTransactionMemo("topic-tools").
		SetAdminKey(client.GetOperatorPublicKey()).
		Execute(client)

	if err != nil {
		println(err.Error(), ": error creating topic")
		return
	}

	// Get the receipt
	transactionReceipt, err := transactionResponse.GetReceipt(client)

	if err != nil {
		println(err.Error(), ": error getting topic create receipt")
		return
	}

	// get the topic id from receipt
	topicID := *transactionReceipt.TopicID

	fmt.Printf("topicID: %v\n", topicID)

}
