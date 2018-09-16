package main

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	_ "github.com/networth-app/networth/api/lib/dotenv"
	"github.com/stretchr/testify/assert"
)

var (
	institutionID = "ins_1"
	token         = "access-sandbox-f9a0d88f-622b-4763-98e5-707692762a50"
)

func TestSyncAccounts(t *testing.T) {
	username := "test_sync_account@networth.app"

	err := syncAccounts(username, institutionID, token)
	assert.Equal(t, nil, err)
}

func TestAppendAccount(t *testing.T) {
	username := "test_append_account@networth.app"

	input := []byte(`
	{ "M":
			{
					"account_id": { "S": "Joe" }
			}
	}`)

	var account events.DynamoDBAttributeValue
	json.Unmarshal(input, &account)

	accounts := []events.DynamoDBAttributeValue{
		account,
	}

	err := appendAccount(username, accounts)
	assert.Equal(t, nil, err)
}
