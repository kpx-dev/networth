package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/api/lib"
)

// append token from single institution to the "all" institution sort key
func appendToken(username string, tokenMap map[string]events.DynamoDBAttributeValue) error {
	token := &nwlib.Token{
		ItemID:          tokenMap["item_id"].String(),
		AccessToken:     tokenMap["access_token"].String(),
		InstitutionID:   tokenMap["institution_id"].String(),
		InstitutionName: tokenMap["institution_name"].String(),
	}

	return db.SetToken(username, nwlib.DefaultSortValue, token)
}
