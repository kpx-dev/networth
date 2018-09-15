package main

import (
	"fmt"
	"log"

	"github.com/networth-app/networth/api/lib"
)

func syncNetworth(username string) {
	log.Printf("Syncing networth for %s", username)
	accounts, err := db.GetAccounts(username, nwlib.DefaultSortValue)

	if err != nil {
		log.Println("Problem getting accounts ", err)
	}

	accountCache := make(map[string]string)
	// var payload
	for _, account := range accounts {
		fmt.Println(account)
	}
}
