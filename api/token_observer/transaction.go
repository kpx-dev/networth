package main

import (
	"fmt"
	"log"

	"github.com/networth-app/networth/api/lib"
)

func transactions(username string, accessTokens []string) {
	for _, token := range accessTokens {
		trans, err := plaid.GetTransactions(token, "2018-01-01", "2018-03-01")

		if err != nil {
			log.Println("Problem getting trans ", err)
			return
		}

		msg := fmt.Sprintf("Total trans: %d", trans.TotalTransactions)
		for _, tran := range trans.Transactions {
			fmt.Println(tran.AccountID, tran.Amount, tran.Date, tran.Name)
		}

		log.Println(msg)
		nwlib.Alert(msg)
	}
}
