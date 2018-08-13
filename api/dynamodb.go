package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

var accountTable = aws.String(getEnv("ACCOUNT_TABLE", ""))

// DBClient db client struct
type DBClient struct {
	*dynamodb.DynamoDB
}

// NewDBClient new dynamodb client
func NewDBClient() *DBClient {
	cfg := loadAWSConfig()
	table := dynamodb.New(cfg)

	return &DBClient{table}
}

// GetTokens return tokens from db
func (d DBClient) GetTokens(username string) []string {
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: accountTable,
		Key: map[string]dynamodb.AttributeValue{
			"username": {S: aws.String(fmt.Sprintf("%s:tokens", username))},
		},
	})

	res, err := req.Send()
	if err != nil {
		panic(err.Error())
	}

	tokens := make(map[string]interface{})
	if err := dynamodbattribute.UnmarshalMap(res.Item, &tokens); err != nil {
		panic(err)
	}

	// TODO: might not be a right approach to init kms client here...
	kmsClient := NewKMSClient()

	payload := []string{""}
	for k, v := range tokens {
		if strings.HasPrefix(k, "ins_") {
			decrypted := kmsClient.Decrypt(v.([]string)[0])
			payload = append(payload, decrypted)
		}
	}

	return payload
}

// GetAccounts return accounts from db
func (d DBClient) GetAccounts(username string) map[string]interface{} {
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: accountTable,
		Key: map[string]dynamodb.AttributeValue{
			"username": {S: aws.String(fmt.Sprintf("%s:accounts", username))},
		},
	})

	res, err := req.Send()
	if err != nil {
		panic(err.Error())
	}

	account := make(map[string]interface{})
	if err := dynamodbattribute.UnmarshalMap(res.Item, &account); err != nil {
		panic(err)
	}

	return account
}

// UpsertAccounts update or insert accounts to db
// func (d DBClient) UpsertAccounts(username string, account Account) {
// 	dynoData, err := dynamodbattribute.Marshal(account)

// 	if err != nil {
// 		panic(err)
// 	}

// 	column := fmt.Sprintf("%s:%s", account.Type, account.Mask)

// 	req := d.UpdateItemRequest(&dynamodb.UpdateItemInput{
// 		Key: map[string]dynamodb.AttributeValue{
// 			"username": {S: aws.String(fmt.Sprintf("%s:accounts", username))},
// 		},
// 		TableName: accountTable,
// 		ExpressionAttributeNames: map[string]string{
// 			"#column": column,
// 		},
// 		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
// 			":column": *dynoData,
// 		},
// 		UpdateExpression: aws.String("SET #column = :column"),
// 	})

// 	req.Send()
// }

// UpdateNetworth update latest networth amount
func (d DBClient) UpdateNetworth(username string, networth float64) {
	req := d.PutItemRequest(&dynamodb.PutItemInput{
		Item: map[string]dynamodb.AttributeValue{
			"username": {S: aws.String(fmt.Sprintf("%s:networth", username))},
			"networth": {N: aws.String(strconv.FormatFloat(networth, 'f', 2, 64))},
		},
		TableName: accountTable,
	})

	req.Send()
}
