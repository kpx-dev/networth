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
	apiHost  = nwlib.GetEnv("API_HOST", ":8000")
)

func main() {
	plaidClient := nwlib.NewPlaidClient()
	dbClient := NewDBClient()

	api := &NetworthAPI{
		db:     dbClient,
		router: mux.NewRouter(),
		plaid:  plaidClient,
	}
	api.Start(apiHost)
}
