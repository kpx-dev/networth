package nwlib

import (
	"net/http"

	"github.com/plaid/plaid-go/plaid"
)

// PlaidClient plaid client
type PlaidClient struct {
	Client  *plaid.Client
	Account *plaid.Account
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

	return &PlaidClient{Client: client}
}
