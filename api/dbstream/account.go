package main

import (
	"log"
)

func syncAccounts(username string, itemID string, token string) error {
	log.Printf("Syncing accounts for username: %s \n", username)
	accounts, err := plaidClient.GetAccounts(token)

	if err != nil {
		log.Printf("Problem getting accounts for username: %s\n", username)
		return err
	}

	for _, account := range accounts.Accounts {
		if err := db.SetAccount(username, itemID, &account); err != nil {
			log.Printf("Problem saving account to db for username: %s, itemID: %s\n", username, itemID)
			return err
		}
	}

	return nil
}
