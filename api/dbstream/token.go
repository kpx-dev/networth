package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func handleInsertModifyToken(username string, itemID string, record events.DynamoDBEventRecord) error {
	token := record.Change.NewImage

	accessToken, err := kms.Decrypt(token["access_token"].String())
	if err != nil {
		log.Println("Problem decoding access_token ", err)
		return err
	}

	// TODO: make these into gorutines / wait group workers:
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
	if err := syncAccounts(username, itemID, accessToken); err != nil {
		log.Println("Problem syncing accounts ", err)
		return err
	}

	// if err := syncTransactions(username, accessToken); err != nil {
	// 	log.Println("Problem syncing transactions ", err)
	//  return err
	// }

	// if err := syncNetworth(username); err != nil {
	// 	log.Println("Problem syncing networth ", err)
	//  return err
	// }

	return nil
}
