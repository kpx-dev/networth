package main

import (
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
	"github.com/stretchr/testify/assert"
)

var (
	username     = "test@networth.app"
	itemID       = "1"
	invalidToken = "access-sandbox-f9a0d88f-622b-4763-98e5-707692762a50"
)

func TestSyncAccounts(t *testing.T) {
	// using invalid token
	err := syncAccounts(username, itemID, invalidToken)
	assert.Equal(t, err != nil, true)
}
