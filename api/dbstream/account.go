package main

import (
	"fmt"
	"log"
)

func accounts(username string, accessTokens []string) error {
	for _, token := range accessTokens {
		accounts, err := plaid.GetAccounts(token)

		if err != nil {
			log.Println("Problem getting accounts ", err)
			return err
		}

		// msg := fmt.Sprintf("Total trans: %d", trans.TotalTransactions)
		for _, account := range accounts.Accounts {
			fmt.Println(account)
		}
	}

	return nil
}
