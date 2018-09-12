package main

import (
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

var (
	institutionID = "ins_1"
)

func TestAccounts(t *testing.T) {
	username := "test_sync_account@networth.app"
	token := "access-sandbox-f9a0d88f-622b-4763-98e5-707692762a50"

	if err := syncAccounts(username, institutionID, token); err != nil {
		t.Error("Failed to parse accounts", err)
	}
}
