package nwlib

import (
	"net/http"

	"github.com/plaid/plaid-go/plaid"
)

// Networth holds networth info
type Networth struct {
	Networth float64 `json:"networth"`
}

// Tokens holds the structure multiple tokens
type Tokens struct {
	Tokens []*Token `json:"tokens"`
}

// Token holds the structure single token
type Token struct {
	ItemID      string `json:"item_id"`
	AccessToken string `json:"access_token"`
	// AccountID       string   `json:"account_id"`
	InstitutionID   string `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
	// Accounts        []string `json:"accounts"`
}

// Transaction struct
type Transaction struct {
	plaid.Transaction
}

// Account wrapper struct for plaid.Account
type Account struct {
	plaid.Account
}

// Accounts hols the structure for multiple plaid account
type Accounts struct {
	Accounts []*Account `json:"accounts"`
}

// PlaidClient plaid client
type PlaidClient struct {
	*plaid.Client
}

// NewPlaidClient new Plaid client
func NewPlaidClient(clientID string, secret string, publicKey string, environment string) *PlaidClient {
	plaidHost := plaid.Sandbox

	switch environment {
	case "sandbox":
		plaidHost = plaid.Sandbox
		break
	case "dev":
		plaidHost = plaid.Development
		break
	case "prod":
		plaidHost = plaid.Production
		break
	default:
		plaidHost = plaid.Sandbox
	}

	clientOptions := plaid.ClientOptions{
		ClientID:    clientID,
		Secret:      secret,
		PublicKey:   publicKey,
		Environment: plaidHost,
		HTTPClient:  &http.Client{},
	}
	client, _ := plaid.NewClient(clientOptions)

	return &PlaidClient{client}
}
