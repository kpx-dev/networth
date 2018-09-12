package main

import (
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/api/lib"
	"github.com/plaid/plaid-go/plaid"
)

func syncAccounts(username string, institutionID string, token string) error {
	accounts, err := plaidClient.GetAccounts(token)

	if err != nil {
		log.Println("syncAccounts() Problem getting accounts ", err)
		return err
	}

	for _, account := range accounts.Accounts {
		db.SetAccount(username, institutionID, &account)
	}

	return nil
}

// append token from single institution to the "all" institution sort key
func appendAccount(username string, accounts []events.DynamoDBAttributeValue) error {
	for _, account := range accounts {
		accountMap := account.Map()

		newAccount := &plaid.Account{
			AccountID: accountMap["account_id"].String(),
		}

		db.SetAccount(username, nwlib.DefaultSortValue, newAccount)
	}

	return nil
}
