package nwlib

import (
	"net/http"

	"github.com/plaid/plaid-go/plaid"
)

// PlaidClient plaid client
type PlaidClient struct {
	*plaid.Client
}

// NewPlaidClient new Plaid client
func NewPlaidClient() *PlaidClient {
	plaidClientID := GetEnv("PLAID_CLIENT_ID")
	plaidSecret := GetEnv("PLAID_SECRET")
	plaidPublicKey := GetEnv("PLAID_PUBLIC_KEY")
	plaidEnv := GetEnv("PLAID_ENV", "sandbox")

	plaidHost := plaid.Sandbox

	switch plaidEnv {
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
		ClientID:    plaidClientID,
		Secret:      plaidSecret,
		PublicKey:   plaidPublicKey,
		Environment: plaidHost,
		HTTPClient:  &http.Client{},
	}
	client, _ := plaid.NewClient(clientOptions)

	return &PlaidClient{client}
}

// // Account account
// type Account struct {
// 	AccountID string `json:"account_id"`
// 	Balances  struct {
// 		Available float64 `json:"available"`
// 		Current   float64 `json:"current"`
// 		// Limit   float64 `json:"limit"`
// 	} `json:"balances"`
// 	Mask         string `json:"mask"`
// 	Name         string `json:"name"`
// 	OfficialName string `json:"official_name"`
// 	Subtype      string `json:"subtype"`
// 	Type         string `json:"type"`
// }

// // Response response
// type Response struct {
// 	// Normal response fields
// 	RequestID string `json:"request_id"`
// 	// Item	struct {} `json:item`
// 	AccessToken      string        `json:"access_token"`
// 	AccountID        string        `json:"account_id"`
// 	Accounts         []Account     `json:"accounts"`
// 	BankAccountToken string        `json:"stripe_bank_account_token"`
// 	MFA              string        `json:"mfa"`
// 	Transactions     []Transaction `json:"transactions"`
// }

// // Transaction transaction
// type Transaction struct {
// 	ID        string  `json:"_id"`
// 	AccountID string  `json:"_account"`
// 	Amount    float64 `json:"amount"`
// 	Date      string  `json:"date"`
// 	Name      string  `json:"name"`
// 	Meta      struct {
// 		AccountOwner string `json:"account_owner"`

// 		Location struct {
// 			Address string `json:"address"`
// 			City    string `json:"city"`

// 			Coordinates struct {
// 				Lat float64 `json:"lat"`
// 				Lon float64 `json:"lon"`
// 			} `json:"coordinates"`

// 			State string `json:"state"`
// 			Zip   string `json:"zip"`
// 		} `json:"location"`
// 	} `json:"meta"`

// 	Pending bool `json:"pending"`

// 	Type struct {
// 		Primary string `json:"primary"`
// 	} `json:"type"`

// 	Category   []string `json:"category"`
// 	CategoryID string   `json:"category_id"`

// 	Score struct {
// 		Location struct {
// 			Address float64 `json:"address"`
// 			City    float64 `json:"city"`
// 			State   float64 `json:"state"`
// 			Zip     float64 `json:"zip"`
// 		}
// 		Name float64 `json:"name"`
// 	} `json:"score"`
// }

// var baseURL = map[string]string{
// 	"sandbox":     "https://sandbox.plaid.com",
// 	"development": "https://development.plaid.com",
// 	"production":  "https://production.plaid.com",
// }

// // NewPlaidClient new Plaid client
// func NewPlaidClient(env string, clientID string, secret string, accessToken string) *PlaidClient {
// 	return &PlaidClient{env, clientID, secret, accessToken}
// }

// func (c PlaidClient) get(path string) *http.Response {
// 	url := fmt.Sprintf("%s/%s", baseURL[c.env], path)
// 	body := map[string]string{
// 		"client_id":    c.clientID,
// 		"secret":       c.secret,
// 		"access_token": c.AccessToken,
// 	}

// 	jsonBody, _ := json.Marshal(body)

// 	res, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))

// 	if err != nil {
// 		fmt.Println("Plaid get error ", err)
// 	}

// 	return res
// }

// // GetAccounts get accounts
// func (c PlaidClient) GetAccounts() Response {
// 	res := c.get("accounts/get")
// 	defer res.Body.Close()

// 	var body Response
// 	err := json.NewDecoder(res.Body).Decode(&body)
// 	if err != nil {
// 		fmt.Println("Plaid decode error ", err)
// 	}

// 	return body
// }
