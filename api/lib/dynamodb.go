package nwlib

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/dynamodbattribute"
	"github.com/plaid/plaid-go/plaid"
)

// Networth holds networth info
type Networth struct {
}

// Tokens holds the structure multiple tokens
type Tokens struct {
	Tokens []*Token `json:"tokens"`
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

// Account wrapper struct for plaid.Account
type Account struct {
	plaid.Account
}

// Accounts hols the structure for multiple plaid account
type Accounts struct {
	Accounts []*Account `json:"accounts"`
}

var (
	networthTable = GetEnv("NETWORTH_TABLE")
	// DefaultSortValue default sort key value
	DefaultSortValue = "all"
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
func (d DynamoDBClient) SetNetworth(username string, networth float64, assets float64, liabilities float64) error {
	now := time.Now().UTC()
	timestamp := now.Format(time.RFC3339)
	networthStr := aws.String(strconv.FormatFloat(networth, 'f', -1, 64))
	assetsStr := aws.String(strconv.FormatFloat(assets, 'f', -1, 64))
	liabilitiesStr := aws.String(strconv.FormatFloat(liabilities, 'f', -1, 64))
	key := fmt.Sprintf("%s:networth", username)

	fmt.Println("SetNetworth for ", username, networth, assets, liabilities)
	req := d.BatchWriteItemRequest(&dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]dynamodb.WriteRequest{
			networthTable: {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"key":         {S: aws.String(key)},
							"sort":        {S: aws.String(timestamp)},
							"networth":    {N: networthStr},
							"assets":      {N: assetsStr},
							"liabilities": {N: liabilitiesStr},
						},
					},
				},
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"key":         {S: aws.String(key)},
							"sort":        {S: aws.String(DefaultSortValue)},
							"networth":    {N: networthStr},
							"assets":      {N: assetsStr},
							"liabilities": {N: liabilitiesStr},
						},
					},
				},
			},
		},
	})

	res, err := req.Send()
	fmt.Printf("Set networth res %+v\n", res)

	return err
}

// GetToken return tokens from db
func (d DynamoDBClient) GetToken(username string, institutionID string) *Tokens {
	dbTokens := &Tokens{}
	key := fmt.Sprintf("%s:token", username)
	sort := DefaultSortValue
	if len(institutionID) > 0 {
		sort = institutionID
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

		return dbTokens
	}

	if err := dynamodbattribute.UnmarshalMap(res.Item, &dbTokens); err != nil {
		log.Println("Problem converting token data from db ", err)

		return dbTokens
	}

	return dbTokens
}

// SetToken save token to db
func (d DynamoDBClient) SetToken(username string, institutionID string, token *Token) error {
	tokenList := [1]*Token{token}
	tokenAttr, err := dynamodbattribute.Marshal(tokenList)
	if err != nil {
		fmt.Println("Problem marshalling token struct into dyno format", err)
		return err
	}

	req := d.UpdateItemRequest(&dynamodb.UpdateItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"key":  {S: aws.String(fmt.Sprintf("%s:token", username))},
			"sort": {S: aws.String(institutionID)},
		},
		TableName: aws.String(networthTable),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":token":      *tokenAttr,
			":emptyToken": {L: []dynamodb.AttributeValue{}},
		},
		UpdateExpression: aws.String("SET tokens = list_append(if_not_exists(tokens, :emptyToken), :token)"),
	})

	if _, err := req.Send(); err != nil {
		log.Println("Problem SetToken ", err)
		return err
	}

	return nil
}

// SetAccount save account to db
func (d DynamoDBClient) SetAccount(username string, institutionID string, account *plaid.Account) error {
	accounts := [1]*plaid.Account{account}

	accountAttr, err := dynamodbattribute.Marshal(accounts)
	if err != nil {
		fmt.Println("Problem marshalling account struct into dyno format", err)
		return err
	}

	req := d.UpdateItemRequest(&dynamodb.UpdateItemInput{
		Key: map[string]dynamodb.AttributeValue{
			"key":  {S: aws.String(fmt.Sprintf("%s:account", username))},
			"sort": {S: aws.String(institutionID)},
		},
		TableName: aws.String(networthTable),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":account":      *accountAttr,
			":emptyAccount": {L: []dynamodb.AttributeValue{}},
		},
		UpdateExpression: aws.String("SET accounts = list_append(if_not_exists(accounts, :emptyAccount), :account)"),
	})

	if _, err := req.Send(); err != nil {
		log.Println("Problem SetAccount ", err)
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
func (d DynamoDBClient) GetAccounts(username string, sort string) (Accounts, error) {
	var accounts Accounts
	key := fmt.Sprintf("%s:account", username)
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(networthTable),
		Key: map[string]dynamodb.AttributeValue{
			"key":  {S: aws.String(key)},
			"sort": {S: aws.String(sort)},
		},
	})

	res, err := req.Send()
	if err != nil {
		return accounts, err
	}

	if err := dynamodbattribute.UnmarshalMap(res.Item, &accounts); err != nil {
		return accounts, err
	}

	return accounts, nil
}

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
