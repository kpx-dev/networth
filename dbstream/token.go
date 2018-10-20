package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/lib"
)

func handleInsertModifyToken(username string, itemID string, record events.DynamoDBEventRecord) error {
	token := record.Change.NewImage

	accessToken, err := kms.Decrypt(token["access_token"].String())
	if err != nil {
		log.Printf("Problem decoding access token: %+v\n", err)
		return err
	}

	// TODO: make these into gorutines / wait group workers:
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
	if err := nwlib.SyncAccounts(plaidClient, db, username, itemID, accessToken); err != nil {
		log.Printf("Problem syncing accounts: %+v\n", err)
		return err
	}

	if err := nwlib.SyncNetworth(db, username); err != nil {
		log.Printf("Problem syncing networth: %+v\n", err)
		return err
	}

	if err := syncTransactions(username, accessToken); err != nil {
		log.Printf("Problem syncing transactions: %+v", err)
		return err
	}

	return nil
}
