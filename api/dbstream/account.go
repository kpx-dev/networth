package main

import (
	"fmt"
	"log"
)

func syncAccounts(username string, token string) error {
	accounts, err := plaid.GetAccounts(token)

	if err != nil {
		log.Println("syncAccounts() Problem getting accounts ", err)
		return err
	}

	for _, account := range accounts.Accounts {
		fmt.Printf("sync account %+v", account)
	}

	return nil
}
