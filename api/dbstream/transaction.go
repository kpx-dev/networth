package main

import (
	"fmt"
	"log"
	"time"
)

func transactions(username string, accessTokens []string) error {
	today := time.Now().UTC()
	todayStr := today.Format("2006-01-02")
	last3Months := today.AddDate(0, -3, 0)
	last3MonthsStr := last3Months.Format("2006-01-02")

	for _, token := range accessTokens {
		trans, err := plaid.GetTransactions(token, last3MonthsStr, todayStr)

		if err != nil {
			log.Println("Problem getting trans ", err)
			return err
		}

		// msg := fmt.Sprintf("Total trans: %d", trans.TotalTransactions)
		for _, account := range trans.Accounts {
			// fmt.Println(tran.AccountID, tran.Amount, tran.Date, tran.Name)
			fmt.Println(account)
		}

		// for _, tran := range trans.Transactions {
		// 	// fmt.Println(tran.AccountID, tran.Amount, tran.Date, tran.Name)
		// 	fmt.Println(tran)
		// }

		// log.Println(msg)
		// nwlib.Alert(msg)
	}

	return nil
}
