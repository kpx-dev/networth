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

	webhookError := newRecord["error"].Map()
	webhookErrorCode := webhookError["error_code"].String()

	switch webhook.WebhookCode {
	case "ERROR":
		if webhookErrorCode == "ITEM_LOGIN_REQUIRED" {
			username, err := db.GetUsernameByItemID(webhook.ItemID)

			if err != nil || username == "" {
				log.Printf("Cannot get username for itemID: %s %+v", webhook.ItemID, err)
			}

			token := &nwlib.Token{
				Sort:     webhook.ItemID,
				Username: username,
				Error:    webhookErrorCode,
			}

			if err := db.UpdateTokenError(token); err != nil {
				log.Printf("Problem updating error value for user: %s\n %+v\n", username, err)
				return err
			}
		}
	default:
		token, err := db.GetTokenByItemID(kms, webhook.ItemID)

		if err != nil {
			log.Printf("Problem getting token for item id: %s \n %+v\n", webhook.ItemID, err)
			return err
		}

		err = sync(token)
	}

	return nil
}
