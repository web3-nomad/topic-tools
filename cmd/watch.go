package main

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
	"os"
	"strings"
	"time"
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

	topicID, err := hedera.TopicIDFromString(os.Getenv("TOPIC_ID"))
	if err != nil {
		println(err.Error(), ": error converting string to topicID")
		return
	}
	println("Watching the topic")
	_, err = hedera.NewTopicMessageQuery().
		// For which topic ID
		SetTopicID(topicID).
		// When to start
		SetStartTime(time.Unix(0, 0)).
		Subscribe(client, func(message hedera.TopicMessage) {
			fmt.Printf("Received message %d\n", message.SequenceNumber)
			fmt.Printf("%s\n", string(message.Contents))
		})

	if err != nil {
		panic(fmt.Sprintf("%v : error subscribing to the topic", err))
		return
	}
	for {
		// Sleep to make sure everything propagates
		time.Sleep(30 * time.Second)
		println(".. still watching ..")
	}

	//println()
}
