package nwlib

import (
	"testing"

	"github.com/stretchr/testify/assert"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

func TestGetToken(t *testing.T) {
	db := NewDynamoDBClient()
	username := "test@networth.app"

	tokens := db.GetToken(username, "")

	assert.Equal(t, len(tokens.Tokens) > 0, true)
	assert.Equal(t, tokens.Tokens[0].InstitutionID, "ins_1")
}
