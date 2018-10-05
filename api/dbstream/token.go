package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func handleInsertModifyToken(username string, sort string, record events.DynamoDBEventRecord) {
	tokens := record.Change.NewImage["tokens"].List()
	newToken := tokens[len(tokens)-1].Map()

	accessToken, err := kms.Decrypt(newToken["access_token"].String())
	if err != nil {
		log.Println("Problem decoding access_token at handleInsertModifyToken")
		return
	}

	// TODO: make these into gorutines / wait group workers:
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
	if err := syncAccounts(username, sort, accessToken); err != nil {
		log.Println("Problem syncing accounts ", err)
	}

	if err := syncTransactions(username, accessToken); err != nil {
		log.Println("Problem syncing transactions ", err)
	}
}
