package main

import (
	"github.com/gorilla/mux"
	"github.com/networth-app/networth/api/lib"
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
	kmsKeyAlias    string
)

func main() {
	accessToken = nwlib.GetEnv("PLAID_ACCESS_TOKEN")
	jwtSecret = getEnv("JWT_SECRET")
	plaidClientID = getEnv("PLAID_CLIENT_ID")
	plaidSecret = getEnv("PLAID_SECRET")
	plaidPublicKey = getEnv("PLAID_PUBLIC_KEY")
	plaidEnv = getEnv("PLAID_ENV", "sandbox")
	kmsKeyAlias = getEnv("KMS_KEY_ALIAS", "alias/networth")
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
