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

	nwlib.PublishSNS(snsARN, "about to sync trans...")
	for _, tran := range trans.Transactions {
		nwlib.PublishSNS(snsARN, tran.ID)
		fmt.Printf("sync transaction %+v\n", tran)

		if err := db.SetTransaction(username, tran); err != nil {
			log.Printf("Problem saving this transaction to db: %+v", tran)
		}
	}

	return nil
}
