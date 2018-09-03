package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/networth-app/networth/api/lib"
)

var (
	tokenTable       = nwlib.GetEnv("TOKEN_TABLE")
	transactionTable = nwlib.GetEnv("TRANSACTION_TABLE")
	historyTable     = nwlib.GetEnv("HISTORY_TABLE")
)

// DynamoDBClient db client struct
type DynamoDBClient struct {
	*dynamodb.DynamoDB
}

// HistoryResp history response
type HistoryResp struct {
	Networth float64 `json:"networth"`
}

// NewDynamoDBClient new dynamodb client
func NewDynamoDBClient() *DynamoDBClient {
	cfg := nwlib.LoadAWSConfig()
	table := dynamodb.New(cfg)

	return &DynamoDBClient{table}
}

// GetNetworth return networth
func (d DynamoDBClient) GetNetworth() float64 {
	today := time.Now().UTC().Format("2006-01-02")

	networth, err := d.Get(historyTable, username, today)

	if err != nil {
		return 0.0
	}

	return networth
}

// SetNetworth value as of today date and current timestamp
func (d DynamoDBClient) SetNetworth(networth float64) error {
	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	timestamp := now.Format(time.RFC3339)
	networthStr := aws.String(strconv.FormatFloat(networth, 'f', -1, 64))

	req := d.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dynamodb.WriteRequest{
			historyTable: {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"email":    {S: aws.String(username)},
							"datetime": {S: aws.String(today)},
							"networth": {N: networthStr},
						},
					},
				},
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"email":    {S: aws.String(username)},
							"datetime": {S: aws.String(timestamp)},
							"networth": {N: networthStr},
						},
					},
				},
			},
		},
	})

	_, err := req.Send()

	return err
}

// GetToken return tokens from db
func (d DynamoDBClient) GetToken(username string) map[string]interface{} {
	dbToken := make(map[string]interface{})

	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(tokenTable),
		Key: map[string]dynamodb.AttributeValue{
			"email": {S: aws.String(username)},
		},
	})

	res, err := req.Send()
	if err != nil {
		log.Println("Problem getting tokens from db ", err)

		return dbToken
	}

	// fmt.Println("got item ", res.Item)
	if err := dynamodbattribute.UnmarshalMap(res.Item, &dbToken); err != nil {
		log.Println("Problem converting token data from db ", err)

		return dbToken
	}

	return dbToken
}

// SetToken save token to db
func (d DynamoDBClient) SetToken(username string, instituionName string, tokenMap *Token) error {
	data, err := dynamodbattribute.Marshal(tokenMap)
	if err != nil {
		fmt.Println("Problem marshalling token map into dyno format", err)
		return err
	}

	req := d.UpdateItemRequest(&dynamodb.UpdateItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"email": {S: aws.String(username)},
		},
		TableName: aws.String(tokenTable),
		ExpressionAttributeNames: map[string]string{
			"#instituion": instituionName,
		},
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":data": *data,
		},
		UpdateExpression: aws.String("SET #instituion = :data"),
	})

	if _, err := req.Send(); err != nil {
		log.Println("Problem saving token to db ", err)
		return err
	}

	return nil
}

// Get item
func (d DynamoDBClient) Get(table string, partitionKey string, sortKey string) (float64, error) {
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]dynamodb.AttributeValue{
			"email":    {S: aws.String(partitionKey)},
			"datetime": {S: aws.String(sortKey)},
		},
	})

	res, err := req.Send()
	if err != nil {
		return 0.0, err
	}

	payload := HistoryResp{}
	if err := dynamodbattribute.UnmarshalMap(res.Item, &payload); err != nil {
		log.Println(err)

		return 0.0, err
	}
	fmt.Println(payload)

	return payload.Networth, nil
}

// Set key / val to db
func (d DynamoDBClient) Set(table string, partitionKey string, sortKey string, valMap map[string]string) error {
	items := map[string]dynamodb.AttributeValue{
		"email": {S: aws.String(partitionKey)},
	}

	if len(sortKey) > 0 {
		items["datetime"] = dynamodb.AttributeValue{S: aws.String(sortKey)}
	}

	for key, val := range valMap {
		fmt.Println("key / val ", key, val)
		items[key] = dynamodb.AttributeValue{S: aws.String(val)}
	}

	fmt.Println(items)

	req := d.PutItemRequest(&dynamodb.PutItemInput{
		Item:      items,
		TableName: aws.String(table),
	})

	res, err := req.Send()

	fmt.Println("Dyno Set() res: ", res)

	return err
}

// func Range() {
// 	req := d.QueryRequest(&dynamodb.QueryInput{
// 		KeyConditionExpression: aws.String("email=:email AND #sort BETWEEN :lastHour AND :now"),
// 		ExpressionAttributeNames: map[string]string{
// 			"#sort": "datetime",
// 		},
// 		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
// 			":email":    {S: aws.String(username)},
// 			":lastHour": {S: aws.String(lastHour)},
// 			":now":      {S: aws.String(now)},
// 		},
// 		Limit:     aws.Int64(1),
// 		TableName: aws.String(getEnv("HISTORY_TABLE")),
// 	})

// 	res, err := req.Send()
// 	if err != nil {
// 		panic(err)
// 	}

// 	if *res.Count > int64(0) {
// 		// nw := make(map[string]interface{})
// 		// if err := dynamodbattribute.UnmarshalMap(res.Items[0], &nw); err != nil {
// 		// 	panic(err)
// 		// }

// 		// fmt.Println(nw)
// 		fmt.Println(res.Items[0]["networth"])
// 		return 1
// 	}

// 	fmt.Println(res)
// }

// GetAccounts return accounts from db
// func (d DynamoDBClient) GetAccounts(username string) map[string]interface{} {
// 	req := d.GetItemRequest(&dynamodb.GetItemInput{
// 		TableName: accountTable,
// 		Key: map[string]dynamodb.AttributeValue{
// 			"username": {S: aws.String(fmt.Sprintf("%s:accounts", username))},
// 		},
// 	})

// 	res, err := req.Send()
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	account := make(map[string]interface{})
// 	if err := dynamodbattribute.UnmarshalMap(res.Item, &account); err != nil {
// 		panic(err)
// 	}

// 	return account
// }

// UpsertAccounts update or insert accounts to db
// func (d DynamoDBClient) UpsertAccounts(username string, account Account) {
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
// func (d DynamoDBClient) UpdateNetworth(username string, networth float64) {
// 	req := d.PutItemRequest(&dynamodb.PutItemInput{
// 		Item: map[string]dynamodb.AttributeValue{
// 			"username": {S: aws.String(fmt.Sprintf("%s:networth", username))},
// 			"networth": {N: aws.String(strconv.FormatFloat(networth, 'f', 2, 64))},
// 		},
// 		TableName: accountTable,
// 	})

// 	req.Send()
// }
