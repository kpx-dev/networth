package main

import (
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

func TestAccounts(t *testing.T) {
	username := "test@networth.app"
	accessTokens := []string{
		"access-sandbox-f9a0d88f-622b-4763-98e5-707692762a50",
		// "access-sandbox-4dcbbdd7-5fa9-480a-8825-fc526169a073",
	}

	if err := accounts(username, accessTokens); err != nil {
		t.Error("Failed to parse accounts", err)
	}
}
