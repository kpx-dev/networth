package main

import (
	"github.com/gorilla/mux"
	"github.com/networth-app/networth/api/lib"
)

// NetworthAPI nw api struct
type NetworthAPI struct {
	db     *DBClient
	router *mux.Router
	plaid  *nwlib.PlaidClient
}

var (
	username = "demo@networth.app"
)

func main() {
	nwlib.LoadDotEnv()
	apiHost := nwlib.GetEnv("API_HOST", ":8000")
	plaidClient := nwlib.NewPlaidClient()
	dbClient := NewDBClient()

	api := &NetworthAPI{
		db:     dbClient,
		router: mux.NewRouter(),
		plaid:  plaidClient,
	}
	api.Start(apiHost)
}
