package main

import (
	"github.com/gorilla/mux"
)

// NetworthAPI nw api struct
type NetworthAPI struct {
	db     *DBClient
	router *mux.Router
	plaid  *PlaidClient
}

var (
	username       = "demo@networth.app"
	accessToken    string
	jwtSecret      string
	plaidEnv       string
	plaidClientID  string
	plaidSecret    string
	plaidPublicKey string
)

func main() {
	accessToken = getEnv("PLAID_ACCESS_TOKEN")
	jwtSecret = getEnv("JWT_SECRET")
	plaidClientID = getEnv("PLAID_CLIENT_ID")
	plaidSecret = getEnv("PLAID_SECRET")
	plaidPublicKey = getEnv("PLAID_PUBLIC_KEY")
	plaidEnv = getEnv("PLAID_ENV", "sandbox")
	apiHost := getEnv("API_HOST", ":8000")

	plaidClient := NewPlaidClient()
	dbClient := NewDBClient()

	api := &NetworthAPI{
		db:     dbClient,
		router: mux.NewRouter(),
		plaid:  plaidClient,
	}
	api.Start(apiHost)
}
