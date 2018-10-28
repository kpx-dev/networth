package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/networth-app/networth/lib"
)

var (
	plaidClientID  = nwlib.GetEnv("PLAID_CLIENT_ID")
	plaidSecret    = nwlib.GetEnv("PLAID_SECRET")
	plaidPublicKey = nwlib.GetEnv("PLAID_PUBLIC_KEY")
	plaidEnv       = nwlib.GetEnv("PLAID_ENV")
	plaidClient    = nwlib.NewPlaidClient(plaidClientID, plaidSecret, plaidPublicKey, plaidEnv)
	kms            = nwlib.NewKMSClient()
	db             = nwlib.NewDynamoDBClient()
	snsARN         = nwlib.GetEnv("SNS_TOPIC_ARN")
)

func handleScheduledEvent(ctx context.Context, e events.CloudWatchEvent) {
	users, err := db.GetAllUsers()
	processedUsers := make(map[string]string)

	if err != nil {
		log.Printf("Problem getting all users: %+v", err)
		return
	}

	for _, user := range users {
		// ignore processed users:
		if _, existed := processedUsers[user.Username]; existed {
			log.Printf("Skipping sync for processed username: %s", user.Username)
			continue
		}

		log.Printf("Sync started for username: %s\n", user.Username)

		tokens, err := db.GetTokens(kms, user.Username)
		if err != nil {
			log.Printf("Problem getting tokens: %+v", err)
		}

		for _, token := range tokens {
			// ignore token with known error:
			if token.Error != "" {
				log.Printf("Skipping sync for token with error: %s", token.Error)
				continue
			}

			if err := nwlib.SyncAccounts(plaidClient, db, &token); err != nil {
				errMsg := fmt.Sprintf("Problem syncing accounts for username: %s, item id: %s\n %+v", user.Username, token.ItemID, err)
				log.Println(errMsg)
				nwlib.PublishSNS(snsARN, errMsg)
				panic(err)
			}
		}

		if err := nwlib.SyncNetworth(db, user.Username); err != nil {
			log.Printf("Problem syncing networth: %+v", err)
		}

		processedUsers[user.Username] = user.Username
	}

	log.Println("Sync done.")
}

func main() {
	lambda.Start(handleScheduledEvent)
}
