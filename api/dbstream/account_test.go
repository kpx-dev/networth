package main

import (
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

var (
	institutionID = "ins_1"
	token         = "access-sandbox-f9a0d88f-622b-4763-98e5-707692762a50"
)

func TestSyncAccounts(t *testing.T) {
	username := "test_sync_account@networth.app"

	if err := syncAccounts(username, institutionID, token); err != nil {
		t.Error("Failed to parse accounts", err)
	}
}

func TestAppendAccount(t *testing.T) {
	username := "test_append_account@networth.app"

	if err := appendAccount(username, account); err != nil {
		t.Error("Failed to parse accounts", err)
	}
}
