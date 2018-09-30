package main

import (
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/api/lib"
	"github.com/plaid/plaid-go/plaid"
)

func handleInsertAccount(username string, sort string, record events.DynamoDBEventRecord) {
	if sort == nwlib.DefaultSortValue {
		syncNetworth(username)
	} else if strings.HasPrefix(sort, "ins_") {
		accounts := record.Change.NewImage["accounts"].List()
		appendAccount(username, accounts)
	}
}

func syncAccounts(username string, institutionID string, token string) error {
	log.Printf("%s - sync accounts, ins %s \n", username, institutionID)
	accounts, err := plaidClient.GetAccounts(token)

	if err != nil {
		log.Println("syncAccounts() Problem getting accounts ", err)
		return err
	}

	for _, account := range accounts.Accounts {
		db.SetAccount(username, institutionID, &account)
	}

	return nil
}

// append token from single institution to the "all" institution sort key
func appendAccount(username string, accounts []events.DynamoDBAttributeValue) error {
	for _, account := range accounts {
		accountMap := account.Map()

		balance := accountMap["balances"].Map()
		avail, _ := balance["available"].Float()
		current, _ := balance["current"].Float()
		limit, _ := balance["limit"].Float()
		unofficialCurrencyCode := ""
		if balance["unofficial_currency_code"].DataType() == events.DataTypeString {
			unofficialCurrencyCode = balance["unofficial_currency_code"].String()
		}

		plaidBalance := &plaid.AccountBalances{
			Available:              avail,
			Current:                current,
			Limit:                  limit,
			ISOCurrencyCode:        balance["iso_currency_code"].String(),
			UnofficialCurrencyCode: unofficialCurrencyCode,
		}

		newAccount := &plaid.Account{
			AccountID:    accountMap["account_id"].String(),
			Balances:     *plaidBalance,
			Mask:         accountMap["mask"].String(),
			Name:         accountMap["name"].String(),
			OfficialName: accountMap["official_name"].String(),
			Subtype:      accountMap["subtype"].String(),
			Type:         accountMap["type"].String(),
		}

		db.SetAccount(username, nwlib.DefaultSortValue, newAccount)
	}

	return nil
}
