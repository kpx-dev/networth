package nwlib

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
)

// Tokens hold the structure multiple tokens
type Tokens struct {
	Tokens []Token `json:"tokens"`
}

// Token holds the structure single token
type Token struct {
	ItemID          string   `json:"item_id"`
	AccessToken     string   `json:"access_token"`
	AccountID       string   `json:"account_id"`
	InstitutionID   string   `json:"institution_id"`
	InstitutionName string   `json:"institution_name"`
	Accounts        []string `json:"accounts"`
}

var (
	networthTable    = GetEnv("NETWORTH_TABLE")
	defaultSortValue = "latest"
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
	cfg := LoadAWSConfig()
	table := dynamodb.New(cfg)

	return &DynamoDBClient{table}
}

// GetNetworth return networth
func (d DynamoDBClient) GetNetworth(username string) float64 {
	today := time.Now().UTC().Format("2006-01-02")

	networth, err := d.Get(username, today)

	if err != nil {
		return 0.0
	}

	return networth
}

// SetNetworth value as of today date and current timestamp
func (d DynamoDBClient) SetNetworth(username string, networth float64) error {
	now := time.Now().UTC()
	// today := now.Format("2006-01-02")
	timestamp := now.Format(time.RFC3339)
	networthStr := aws.String(strconv.FormatFloat(networth, 'f', -1, 64))

	req := d.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dynamodb.WriteRequest{
			networthTable: {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"key":      {S: aws.String(fmt.Sprintf("%s:networth", username))},
							"sort":     {S: aws.String(timestamp)},
							"networth": {N: networthStr},
						},
					},
				},
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"key":      {S: aws.String(username)},
							"sort":     {S: aws.String(defaultSortValue)},
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
func (d DynamoDBClient) GetToken(username string, institution string) map[string]interface{} {
	dbToken := make(map[string]interface{})
	key := fmt.Sprintf("%s:token", username)
	sort := defaultSortValue
	if len(institution) > 0 {
		sort = institution
	}

	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(networthTable),
		Key: map[string]dynamodb.AttributeValue{
			"key":  {S: aws.String(key)},
			"sort": {S: aws.String(sort)},
		},
	})

	res, err := req.Send()
	if err != nil {
		log.Printf("Problem getting tokens from db using sort key %s %v", sort, err)

		return dbToken
	}

	if err := dynamodbattribute.UnmarshalMap(res.Item, &dbToken); err != nil {
		log.Println("Problem converting token data from db ", err)

		return dbToken
	}

	// TODO: return Token struct instead of interface
	return dbToken
}

// SetToken save token to db
func (d DynamoDBClient) SetToken(username string, tokenMap *Token) error {
	data, err := dynamodbattribute.Marshal(tokenMap)
	if err != nil {
		fmt.Println("Problem marshalling token map into dyno format", err)
		return err
	}

	req := d.UpdateItemRequest(&dynamodb.UpdateItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"key":  {S: aws.String(fmt.Sprintf("%s:token", username))},
			"sort": {S: aws.String(defaultSortValue)},
		},
		TableName: aws.String(networthTable),
		// ExpressionAttributeNames: map[string]string{
		// 	"#institution": tokenMap.InstitutionID,
		// },
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":data": *data,
		},
		UpdateExpression: aws.String("SET token = :data"),
	})

	if _, err := req.Send(); err != nil {
		log.Println("Problem saving token to db ", err)
		return err
	}

	return nil
}

// Get item
func (d DynamoDBClient) Get(partitionKey string, sortKey string) (float64, error) {
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(networthTable),
		Key: map[string]dynamodb.AttributeValue{
			"key":  {S: aws.String(partitionKey)},
			"sort": {S: aws.String(sortKey)},
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
		"key": {S: aws.String(partitionKey)},
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
