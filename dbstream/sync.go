package main

import (
	"log"

	"github.com/networth-app/networth/lib"
)

func sync(token nwlib.Token) error {
	// TODO: make these into gorutines / wait group workers:
	// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
	if err := nwlib.SyncAccounts(plaidClient, db, token.Username, token.ItemID, token.AccessToken); err != nil {
		log.Printf("Problem syncing accounts: %+v\n", err)
		return err
	}

	if err := nwlib.SyncNetworth(db, token.Username); err != nil {
		log.Printf("Problem syncing networth: %+v\n", err)
		return err
	}

	if err := syncTransactions(token.Username, token.AccessToken); err != nil {
		log.Printf("Problem syncing transactions: %+v", err)
		return err
	}

	return nil
}
