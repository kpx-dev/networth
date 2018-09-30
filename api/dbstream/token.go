package main

import (
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/networth-app/networth/api/lib"
)

func handleInsertModifyToken(username string, sort string, record events.DynamoDBEventRecord) {
	if strings.HasPrefix(sort, "ins_") {
		tokens := record.Change.NewImage["tokens"].List()
		newToken := tokens[len(tokens)-1].Map()

		if err := appendToken(username, newToken); err != nil {
			log.Println("Problem append single token to global list ", err)
			return
		}

		accessToken, err := kms.Decrypt(newToken["access_token"].String())
		if err != nil {
			log.Println("Problem decoding access_token at handleInsertModifyToken")
			return
		}

		// TODO: make these into gorutines / wait group workers:
		// http://devs.cloudimmunity.com/gotchas-and-common-mistakes-in-go-golang/index.html#gor_app_exit
		// syncTransactions(username, accessToken)
		if err := syncAccounts(username, sort, accessToken); err != nil {
			log.Println("Problem syncing accounts ", err)
		}
	}
}

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
