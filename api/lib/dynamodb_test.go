package nwlib

import (
	"fmt"
	"testing"

	_ "github.com/networth-app/networth/api/lib/dotenv"
)

func TestGetToken(t *testing.T) {
	db := NewDynamoDBClient()
	username := "test@networth.app"

	tokens := db.GetToken(username, "")

	fmt.Println(tokens)
	// t.Error("Failed to parse transactions", err)
}
