package nwlib

import (
	"net/http"

	"github.com/plaid/plaid-go/plaid"
)

// Token holds the structure single token
type Token struct {
	ID              string   `json:"id"`
	Sort            string   `json:"sort"`
	ItemID          string   `json:"item_id"`
	AccessToken     string   `json:"access_token"`
	AccountID       string   `json:"account_id"`
	InstitutionID   string   `json:"institution_id"`
	InstitutionName string   `json:"institution_name"`
	Accounts        []string `json:"accounts"`
	Username        string   `json:"username"`
	Error           string   `json:"error"`
}

// Transaction struct
type Transaction struct {
	plaid.Transaction
}

// WebhookError struct contains error fields from webhook
type WebhookError struct {
	DisplayMessage string `json:"display_message"` //: "The provided credentials were not correct. Please try again.",
	ErrorCode      string `json:"error_code"`      //: "ITEM_LOGIN_REQUIRED",
	ErrorMessage   string `json:"error_message"`   //: "the provided credentials were not correct",
	ErrorType      string `json:"error_type"`      //: "ITEM_ERROR",
	Status         int    `json:"status"`          //: 400
}

// Webhook body
type Webhook struct {
	WebhookType         string       `json:"webhook_type"`         // "TRANSACTIONS",
	WebhookCode         string       `json:"webhook_code"`         // INITIAL_UPDATE HISTORICAL_UPDATE DEFAULT_UPDATE TRANSACTIONS_REMOVED
	ItemID              string       `json:"item_id"`              // "wz666MBjYWTp2PDzzggYhM6oWWmBb",
	Error               WebhookError `json:"error"`                // null,
	NewTransactions     int          `json:"new_transactions"`     //19
	RemovedTransactions []string     `json:"removed_transactions"` // ["yBVBEwrPyJs8GvR77N7QTxnGg6wG74H7dEDN6", "kgygNvAVPzSX9KkddNdWHaVGRVex1MHm3k9no"],
	Username            string       `json:"username"`
}

// Account wrapper struct for plaid.Account
type Account struct {
	plaid.Account
	InstitutionID   string `json:"institution_id"`
	InstitutionName string `json:"institution_name"`
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
	case "dev", "development":
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
