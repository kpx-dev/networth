package main

import (
	"context"
	"fmt"

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

	if err != nil {
		fmt.Println("Problem getting all users", err)
		return
	}

	for _, user := range users {
		fmt.Printf("Sync started for username: %s\n", user.Username)

		tokens, err := db.GetTokens(kms, user.Username)
		if err != nil {
			fmt.Println("Problem getting tokens ", err)
		}

		for _, token := range tokens {
			if err := nwlib.SyncAccounts(plaidClient, db, user.Username, token.ItemID, token.AccessToken); err != nil {
				errMsg := fmt.Sprintf("Problem syncing accounts for username :%s\n %+v", user.Username, err)
				fmt.Println(errMsg)
				nwlib.PublishSNS(snsARN, errMsg)
				panic(err)
			}
		}

		if err := nwlib.SyncNetworth(db, user.Username); err != nil {
			fmt.Println("Problem syncing networth ", err)
		}
	}

	fmt.Println("Sync done.")
}

func main() {
	lambda.Start(handleScheduledEvent)
}
