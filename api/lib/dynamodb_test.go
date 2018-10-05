package nwlib

import (
	"testing"

	"github.com/plaid/plaid-go/plaid"
	"github.com/stretchr/testify/assert"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

var (
	db            = NewDynamoDBClient()
	institutionID = "ins_1"
	testUsername  = "test@networth.app"
)

func TestSetTransaction(t *testing.T) {
	username := "test_set_transaction@networth.app"
	trans := plaid.Transaction{
		ID:        "1",
		AccountID: "1",
		Amount:    1,
	}

	err := db.SetTransaction(username, trans)
	assert.Equal(t, err, nil)
}

func TestSetAccount(t *testing.T) {
	username := "test_set_account@networth.app"
	account := &plaid.Account{
		AccountID: "1",
		Name:      "test",
	}

	// set for specific ins
	err := db.SetAccount(username, institutionID, account)
	assert.Equal(t, err, nil)

	// set for all
	err = db.SetAccount(username, DefaultSortValue, account)
	assert.Equal(t, err, nil)
}

func TestSetToken(t *testing.T) {
	db := NewDynamoDBClient()
	username := "test_set_token@networth.app"
	token := &Token{
		ItemID:        "1",
		AccessToken:   "1",
		InstitutionID: institutionID,
	}

	err := db.SetToken(username, token)
	assert.Equal(t, nil, err)
}

func TestGetToken(t *testing.T) {
	db := NewDynamoDBClient()
	username := "test_set_token@networth.app"
	invalidInstitutionID := "ins_1_invalid"

	// get without ins_id
	tokens := db.GetToken(username, "")
	assert.Equal(t, tokens.Tokens[len(tokens.Tokens)-1].InstitutionID, institutionID)

	// get using ins_id
	tokens = db.GetToken(username, institutionID)
	assert.Equal(t, tokens.Tokens[len(tokens.Tokens)-1].InstitutionID, institutionID)

	// // get using invalid ins_id
	tokens = db.GetToken(username, invalidInstitutionID)
	assert.Equal(t, len(tokens.Tokens) == 0, true)
}

func TestGetAccounts(t *testing.T) {
	db := NewDynamoDBClient()
	usernameNotExist := "test_get_account@networth.app"
	usernameExist := "test_set_account@networth.app"
	invalidInstitutionID := "ins_1_invalid"

	// get using non-exist username
	accounts, err := db.GetAccounts(usernameNotExist, institutionID)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(accounts.Accounts), 0)

	// get using existing username
	accounts, err = db.GetAccounts(usernameExist, institutionID)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(accounts.Accounts) > 0, true)

	// get using invalid ins_id
	accounts, err = db.GetAccounts(usernameExist, invalidInstitutionID)
	assert.Equal(t, err, nil)
	assert.Equal(t, len(accounts.Accounts), 0)

	// get using default sort key
	accounts, err = db.GetAccounts(usernameExist, DefaultSortValue)
	assert.Equal(t, err, nil)
	assert.Equal(t, true, len(accounts.Accounts) > 0)
}
