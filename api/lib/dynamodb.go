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

var (
	dbTable = GetEnv("DB_TABLE")
	// DefaultSortValue default sort key value
	DefaultSortValue = "all"
)

// DynamoDBClient db client struct
type DynamoDBClient struct {
	*dynamodb.DynamoDB
}

// NewDynamoDBClient new dynamodb client
func NewDynamoDBClient() *DynamoDBClient {
	cfg := LoadAWSConfig()
	table := dynamodb.New(cfg)

	return &DynamoDBClient{table}
}

// GetNetworth return networth
func (d DynamoDBClient) GetNetworth(username string) float64 {
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(dbTable),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(fmt.Sprintf("%s:networth", username))},
			"sort": {S: aws.String(DefaultSortValue)},
		},
	})

	res, err := req.Send()
	if err != nil {
		log.Println("Problem getting networth ", err)
		return 0.0
	}

	payload := Networth{}
	if err := dynamodbattribute.UnmarshalMap(res.Item, &payload); err != nil {
		log.Println("Problem converting db to Networth struct ", err)

		return 0.0
	}

	return payload.Networth
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
			dbTable: {
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"id":          {S: aws.String(key)},
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
							"id":          {S: aws.String(key)},
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
		TableName: aws.String(dbTable),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(key)},
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
func (d DynamoDBClient) SetToken(username string, token *Token) error {
	tokenAttr, err := dynamodbattribute.MarshalMap(token)
	if err != nil {
		fmt.Println("Problem marshalling token struct into dyno format", err)
		return err
	}

	dbKey := map[string]dynamodb.AttributeValue{
		"id":   {S: aws.String(fmt.Sprintf("%s:token", username))},
		"sort": {S: aws.String(token.ItemID)},
	}

	for k, v := range dbKey {
		tokenAttr[k] = v
	}

	req := d.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(dbTable),
		Item:      tokenAttr,
	})

	if _, err := req.Send(); err != nil {
		log.Println("Problem SetToken ", err)
		return err
	}

	return nil
}

// SetTransaction save transaction to db
func (d DynamoDBClient) SetTransaction(username string, transaction plaid.Transaction) error {

	transactionAttr, err := dynamodbattribute.MarshalMap(transaction)
	if err != nil {
		fmt.Println("Problem marshalling transaction struct into dyno format", err)
		return err
	}

	dbKey := map[string]dynamodb.AttributeValue{
		"id":             {S: aws.String(fmt.Sprintf("%s:transaction", username))},
		"sort":           {S: aws.String(transaction.ID)},
		"transaction_id": {S: aws.String(transaction.ID)},
	}

	for k, v := range dbKey {
		transactionAttr[k] = v
	}

	req := d.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(dbTable),
		Item:      transactionAttr,
	})

	if _, err := req.Send(); err != nil {
		log.Println("Problem SetTransaction ", err)
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
			"id":   {S: aws.String(fmt.Sprintf("%s:account", username))},
			"sort": {S: aws.String(institutionID)},
		},
		TableName: aws.String(dbTable),
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

// Set key / val to db
func (d DynamoDBClient) Set(table string, partitionKey string, sortKey string, valMap map[string]string) error {
	items := map[string]dynamodb.AttributeValue{
		"id": {S: aws.String(partitionKey)},
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

// GetAccounts return accounts from db
func (d DynamoDBClient) GetAccounts(username string, sort string) (Accounts, error) {
	var accounts Accounts
	key := fmt.Sprintf("%s:account", username)
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(dbTable),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(key)},
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
