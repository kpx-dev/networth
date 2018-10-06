package main

import (
	"log"

	"github.com/networth-app/networth/api/lib"

	"github.com/aws/aws-lambda-go/events"
)

func handleInsertModifyWebhook(record events.DynamoDBEventRecord) error {
	newRecord := record.Change.NewImage
	webhook := nwlib.Webhook{
		ItemID: newRecord["item_id"].String(),
	}

	if err := db.SetWebhook(webhook); err != nil {
		log.Println("Problem saving webhook ", err)
		return err
	}

	return nil
}
