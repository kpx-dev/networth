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
	username := "test_sync_accounts@networth.app"

	err := syncAccounts(username, institutionID, token)
	assert.Equal(t, nil, err)
}

func TestAppendAccount(t *testing.T) {
	username := "test_append_account@networth.app"

	input := []byte(`
	{ "M":
			{
					"account_id": { "S": "123" },
					"mask": {"S": "123"},
					"name": {"S": "Hola bank"},
					"official_name": {"S": "BOA"},
					"subtype": {"S": "subtype"},
					"type": {"S": "type"},
					"balances": {
						"M": {
							"available": {"N": "1"},
							"current": {"N": "1"},
							"limit": {"N": "1"},
							"unofficial_currency_code": {"S": "what"},
							"iso_currency_code": {"S": "usd"}
						}
					}

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
