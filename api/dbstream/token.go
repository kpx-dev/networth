package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/api/lib"
)

// append token from single institution to the "all" institution sort key
func appendToken(key string, tokenMap map[string]events.DynamoDBAttributeValue) error {
	username := strings.Split(key, ":")
	fmt.Println("username is ", username[0])
	fmt.Println("tokenMap[item_id].String() ", tokenMap["item_id"].String())
	fmt.Println("tokenMap[access_token].String() ", tokenMap["access_token"].String())

	token := &nwlib.Token{
		ItemID:          tokenMap["item_id"].String(),
		AccessToken:     tokenMap["access_token"].String(),
		InstitutionID:   tokenMap["institution_id"].String(),
		InstitutionName: tokenMap["institution_name"].String(),
	}

	return db.SetToken(username[0], nwlib.DefaultSortValue, token)
}
