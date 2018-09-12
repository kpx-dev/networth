package main

import (
	"fmt"
	"log"
	"time"
)

func syncTransactions(username string, token string) error {
	endDate := time.Now().UTC()
	endDateStr := endDate.Format("2006-01-02")
	startDate := endDate.AddDate(0, -3, 0)
	startDateStr := startDate.Format("2006-01-02")

	trans, err := plaidClient.GetTransactions(token, startDateStr, endDateStr)

	if err != nil {
		log.Println("syncTransactions() Problem getting trans ", err)
		return err
	}

	fmt.Printf("Total trans: %d", trans.TotalTransactions)
	// for _, account := range trans.Accounts {
	// 	// fmt.Println(tran.AccountID, tran.Amount, tran.Date, tran.Name)
	// 	fmt.Println("sync accounts" account)
	// }

	for _, tran := range trans.Transactions {
		// fmt.Println(tran.AccountID, tran.Amount, tran.Date, tran.Name)
		fmt.Printf("sync transaction %+v", tran)
	}

	return nil
}
