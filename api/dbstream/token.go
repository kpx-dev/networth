package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/api/lib"
)

// func decryptTokens() tokens []string {

// }

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
	// fmt.Printf("token is %+v", token)
	return db.SetToken(username[0], nwlib.DefaultSortValue, token)
	// return nil
}

func tokens(record events.DynamoDBEventRecord) (username string, tokens []string) {
	var email string
	var decryptedAccessTokens []string

	for name, value := range record.Change.NewImage {
		if name == "email" {
			email = value.String()
		}

		if value.DataType() == events.DataTypeMap {
			val := value.Map()
			tokens := val["access_tokens"].List()

			for _, token := range tokens {
				decryptedToken := kms.Decrypt(token.String())
				decryptedAccessTokens = append(decryptedAccessTokens, decryptedToken)
			}
		}
	}

	return email, decryptedAccessTokens
}
