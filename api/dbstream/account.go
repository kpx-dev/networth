package main

import (
	"log"
)

func syncAccounts(username string, institutionID string, token string) error {
	accounts, err := plaid.GetAccounts(token)

	if err != nil {
		log.Println("syncAccounts() Problem getting accounts ", err)
		return err
	}

	for _, account := range accounts.Accounts {
		go db.SetAccount(username, institutionID, &account)
	}

	return nil
}

// append token from single institution to the "all" institution sort key
// func appendAccount(username string, accountMap map[string]events.DynamoDBAttributeValue) error {
// 	account := &plaid.Account{
// 		ItemID:          tokenMap["item_id"].String(),
// 		AccessToken:     tokenMap["access_token"].String(),
// 		InstitutionID:   tokenMap["institution_id"].String(),
// 		InstitutionName: tokenMap["institution_name"].String(),
// 	}

// 	return db.SetToken(username, nwlib.DefaultSortValue, token)
// }
