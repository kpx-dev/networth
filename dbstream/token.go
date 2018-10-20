package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/lib"
)

func handleInsertModifyToken(record events.DynamoDBEventRecord) error {
	newToken := record.Change.NewImage
	token := nwlib.Token{
		Username:    newToken["username"].String(),
		ItemID:      newToken["item_id"].String(),
		AccessToken: newToken["access_token"].String(),
	}

	accessToken, err := kms.Decrypt(token.AccessToken)
	if err != nil {
		log.Printf("Problem decoding access token: %+v\n", err)
		return err
	}
	token.AccessToken = accessToken

	err = sync(token)

	return err
}
