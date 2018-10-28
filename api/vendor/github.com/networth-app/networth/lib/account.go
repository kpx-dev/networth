package nwlib

import (
	"log"
)

// SyncAccounts - get latest accounts from Plaid
func SyncAccounts(plaidClient *PlaidClient, db *DynamoDBClient, token *Token) error {
	log.Printf("Syncing accounts for username: %s \n", token.Username)
	accounts, err := plaidClient.GetAccounts(token.AccessToken)

	if err != nil {
		log.Printf("Problem getting accounts for username: %s\n", token.Username)
		return err
	}

	for _, account := range accounts.Accounts {
		if err := db.SetAccount(token, &account); err != nil {
			log.Printf("Problem saving account to db for username: %s, itemID: %s\n", token.Username, token.ItemID)
			return err
		}
	}

	return nil
}
