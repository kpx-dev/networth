package main

import (
	"fmt"
	"log"

	"github.com/networth-app/networth/api/lib"
)

func isAsset(account *nwlib.Account) bool {
	switch account.Type {
	case "brokerage", "depository":
		return true
	case "loan", "mortgage":
		return false
	}

	switch account.Subtype {
	case "credit card", "line of credit":
		return false
	case "paypal", "403B", "cash management", "cd", "hsa", "keogh", "money market", "mutual fund", "prepaid", "recurring", "rewards", "safe deposit", "sarsep":
		return true
	}

	return false
}

func syncNetworth(username string) {
	asset := 0.0
	liability := 0.0
	accountCache := make(map[string]bool)
	accounts, err := db.GetAccounts(username, nwlib.DefaultSortValue)

	if err != nil {
		log.Println("Problem getting accounts ", err)
		return
	}

	for _, account := range accounts.Accounts {
		fmt.Printf("debug account %+v", account)

		if _, ok := accountCache[account.AccountID]; !ok {
			if isAsset(account) {
				asset = asset + account.Balances.Current
			} else {
				liability = liability + account.Balances.Current
			}
			accountCache[account.AccountID] = true
		}
	}

	networth := asset - liability
	msg := fmt.Sprintf("Networth for %s is %f\n", username, networth)
	log.Printf(msg)
	if err := db.SetNetworth(username, networth); err != nil {
		log.Println("Problem setting networth ", err)
		return
	}

	nwlib.PublishSNS(snsARN, msg)
}
