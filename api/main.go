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
	awsRegion      string
)

func main() {
	awsRegion = nwlib.GetEnv("AWS_REGION", "us-east-1")
	accessToken = nwlib.GetEnv("PLAID_ACCESS_TOKEN")
	jwtSecret = nwlib.GetEnv("JWT_SECRET")
	plaidClientID = nwlib.GetEnv("PLAID_CLIENT_ID")
	plaidSecret = nwlib.GetEnv("PLAID_SECRET")
	plaidPublicKey = nwlib.GetEnv("PLAID_PUBLIC_KEY")
	plaidEnv = nwlib.GetEnv("PLAID_ENV", "sandbox")
	kmsKeyAlias = nwlib.GetEnv("KMS_KEY_ALIAS", "alias/networth")
	apiHost := nwlib.GetEnv("API_HOST", ":8000")

	plaidClient := NewPlaidClient()
	dbClient := NewDBClient()

	api := &NetworthAPI{
		db:     dbClient,
		router: mux.NewRouter(),
		plaid:  plaidClient,
	}
	api.Start(apiHost)
}
