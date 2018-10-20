package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/lib"
)

func handleInsertModifyWebhook(record events.DynamoDBEventRecord) error {
	newRecord := record.Change.NewImage
	webhook := nwlib.Webhook{
		ItemID:      newRecord["item_id"].String(),
		WebhookType: newRecord["webhook_type"].String(),
		WebhookCode: newRecord["webhook_code"].String(),
	}

	token, err := db.GetTokenByItemID(kms, webhook.ItemID)

	if err != nil {
		log.Printf("Problem getting token for item id: %s \n %+v\n", webhook.ItemID, err)
	}

	err = sync(token)

	return nil
}
