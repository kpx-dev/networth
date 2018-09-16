package main

import (
	"fmt"
	"log"

	"github.com/networth-app/networth/api/lib"
)

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
			fmt.Println("not yet cached...")
			asset = asset + account.Balances.Current
			accountCache[account.AccountID] = true
		}
	}

	networth := asset - liability
	log.Printf("Networth for %s is %f\n", username, networth)
	if err := db.SetNetworth(username, networth); err != nil {
		log.Println("Problem setting networth ", err)
	}
}
