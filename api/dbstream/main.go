package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/networth-app/networth/api/lib"
)

var (
	plaidClientID  = nwlib.GetEnv("PLAID_CLIENT_ID")
	plaidSecret    = nwlib.GetEnv("PLAID_SECRET")
	plaidPublicKey = nwlib.GetEnv("PLAID_PUBLIC_KEY")
	plaidEnv       = nwlib.GetEnv("PLAID_ENV", "sandbox")
	plaidClient    = nwlib.NewPlaidClient(plaidClientID, plaidSecret, plaidPublicKey, plaidEnv)
	kms            = nwlib.NewKMSClient()
	db             = nwlib.NewDynamoDBClient()
	snsARN         = nwlib.GetEnv("SNS_TOPIC_ARN")
	slackURL       = nwlib.GetEnv("SLACK_WEBHOOK_URL")
	wg             sync.WaitGroup
)

func extractCompositeKeys(record events.DynamoDBEventRecord) (string, string, string) {
	key := record.Change.Keys["id"].String()
	username := strings.Split(key, ":")[0]
	sort := record.Change.Keys["sort"].String()

	return key, username, sort
}

func handleNewToken(username string, sort string, record events.DynamoDBEventRecord) {
	tokens := record.Change.NewImage["tokens"].List()
	newToken := tokens[len(tokens)-1].Map()

	if len(record.Change.OldImage) > 0 {
		appendToken(username, newToken)
	}

	accessToken, err := kms.Decrypt(newToken["access_token"].String())

	if err != nil {
		log.Println("Problem decoding access_token")
		return
	}

	// TODO: make these into gorutines / wait group workers:
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
	// syncTransactions(username, accessToken)
	syncAccounts(username, sort, accessToken)
}

func handleNewAccount(username string, sort string, record events.DynamoDBEventRecord) {
	if sort == nwlib.DefaultSortValue {
		syncNetworth(username)
	} else if strings.HasPrefix(sort, "ins_") && len(record.Change.OldImage) > 0 {
		// each user has 2 keys for account: all, ins_XXX
		// TODO: [bug] somehow have duplicate accounts in all key
		accounts := record.Change.NewImage["accounts"].List()
		appendAccount(username, accounts)
	}
}

// TODO: https://github.com/aws/aws-lambda-go/issues/58
func handleDynamoDBStream(ctx context.Context, e events.DynamoDBEvent) {
	fmt.Printf("handleDynamoDBStream ...")

	for _, record := range e.Records {
		switch record.EventName {
		case "MODIFY":
			fmt.Printf("Modify record %+v", record)
			break
		case "INSERT":
			key, username, sort := extractCompositeKeys(record)

			// each user have 2 sort keys for token: all, ins_XXX
			if strings.HasSuffix(key, ":token") && strings.HasPrefix(sort, "ins_") {
				handleNewToken(username, sort, record)
			} else if strings.HasSuffix(key, ":account") {
				handleNewAccount(username, sort, record)
			}

			break
		default:
			log.Printf("DynamoDB stream unknown event %s %+v", record.EventName, record)
		}
	}
}

func main() {
	lambda.Start(handleDynamoDBStream)
}
