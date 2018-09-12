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
)

func TestSetAccount(t *testing.T) {
	username := "test_set_account@networth.app"
	account := &plaid.Account{
		AccountID: "1",
		Name:      "test",
	}

	// set for specific ins
	if err := db.SetAccount(username, institutionID, account); err != nil {
		t.Errorf("Cannot set account %v", err)
	}
}

func TestSetToken(t *testing.T) {
	db := NewDynamoDBClient()
	username := "test_set@networth.app"
	institutionID := "ins_1"
	token := &Token{
		AccessToken: "1",
	}

	// set for specific ins
	if err := db.SetToken(username, institutionID, token); err != nil {
		t.Errorf("Cannot set token %v", err)
	}
}

func TestGetToken(t *testing.T) {
	db := NewDynamoDBClient()
	username := "test@networth.app"
	institutionID := "ins_1"
	invalidInstitutionID := "ins_1_invalid"

	// get without ins_id
	tokens := db.GetToken(username, "")
	assert.Equal(t, tokens.Tokens[0].InstitutionID, institutionID)

	// get using ins_id
	tokens = db.GetToken(username, institutionID)
	assert.Equal(t, tokens.Tokens[0].InstitutionID, institutionID)

	// get using invalid ins_id
	tokens = db.GetToken(username, invalidInstitutionID)
	assert.Equal(t, len(tokens.Tokens) == 0, true)
}
