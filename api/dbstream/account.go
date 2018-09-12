package main

import (
	"fmt"
	"log"

	"github.com/networth-app/networth/api/lib"
)

func syncAccounts(plaidClient *nwlib.PlaidClient, username string, token string) error {
	accounts, err := plaidClient.GetAccounts(token)

	if err != nil {
		log.Println("syncAccounts() Problem getting accounts ", err)
		return err
	}

	for _, account := range accounts.Accounts {
		fmt.Printf("sync account %+v", account)
	}

	return nil
}
