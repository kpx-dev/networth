package main

import (
	"fmt"
	"log"
	"time"

	"github.com/networth-app/networth/api/lib"
)

// sync last 12 months
func syncTransactions(username string, token string) error {
	endDate := time.Now().UTC()
	endDateStr := endDate.Format("2006-01-02")
	startDate := endDate.AddDate(0, -12, 0)
	startDateStr := startDate.Format("2006-01-02")

	trans, err := plaidClient.GetTransactions(token, startDateStr, endDateStr)

	if err != nil {
		log.Println("syncTransactions() Problem getting trans ", err)
		return err
	}

	nwlib.PublishSNS(snsARN, fmt.Sprintf("trans %+v", trans))
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
