package main

import (
	"fmt"
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"

	"github.com/networth-app/networth/api/lib"
)

var (
	testPlaidClientID  = nwlib.GetEnv("PLAID_CLIENT_ID")
	testPlaidSecret    = nwlib.GetEnv("PLAID_SECRET")
	testPlaidPublicKey = nwlib.GetEnv("PLAID_PUBLIC_KEY")
	testPlaidEnv       = nwlib.GetEnv("PLAID_ENV", "sandbox")
	testPlaidClient    = nwlib.NewPlaidClient(testPlaidClientID, testPlaidSecret, testPlaidPublicKey, testPlaidEnv)
)

func TestAccounts(t *testing.T) {
	fmt.Println("testPlaidClientID", testPlaidClientID)

	username := "test@networth.app"
	token := "access-sandbox-f9a0d88f-622b-4763-98e5-707692762a50"

	if err := syncAccounts(testPlaidClient, username, token); err != nil {
		t.Error("Failed to parse accounts", err)
	}
}
