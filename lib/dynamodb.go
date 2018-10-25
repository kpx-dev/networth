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
)

// DynamoDBClient db client struct
type DynamoDBClient struct {
	*dynamodb.DynamoDB
}

// NewDynamoDBClient new dynamodb client
func NewDynamoDBClient() *DynamoDBClient {
	// TODO: be able to accept custom endpoint (for testing)
	cfg := LoadAWSConfig()
	table := dynamodb.New(cfg)

	return &DynamoDBClient{table}
}

// GetNetworth return Networth struct
func (d DynamoDBClient) GetNetworth(username string) (Networth, error) {
	var networth Networth
	req := d.GetItemRequest(&dynamodb.GetItemInput{
		TableName: aws.String(dbTable),
		Key: map[string]dynamodb.AttributeValue{
			"id":   {S: aws.String(fmt.Sprintf("%s:networth", username))},
			"sort": {S: aws.String("latest")},
		},
	})

	res, err := req.Send()
	if err != nil {
		log.Printf("Problem getting networth: %+v\n", err)
		return networth, err
	}

	payload := Networth{}
	if err := dynamodbattribute.UnmarshalMap(res.Item, &payload); err != nil {
		log.Printf("Problem converting db to Networth struct: %+v\n", err)

		return networth, err
	}

	return payload, nil
}

// GetNetworthByDateRange return net worth based on date
func (d DynamoDBClient) GetNetworthByDateRange(username string, startDate string, endDate string) ([]Networth, error) {
	var networth []Networth
	req := d.QueryRequest(&dynamodb.QueryInput{
		TableName:              aws.String(dbTable),
		KeyConditionExpression: aws.String("id = :id AND sort BETWEEN :startDate AND :endDate"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":id":        {S: aws.String(fmt.Sprintf("%s:networth", username))},
			":startDate": {S: aws.String(startDate)},
			":endDate":   {S: aws.String(endDate)},
		},
	})

	res, err := req.Send()
	if err != nil {
		log.Printf("Problem getting networth by date range: %s - %s %+v", startDate, endDate, err)
		return networth, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &networth); err != nil {
		return networth, err
	}

	return networth, nil
}

// SetNetworth value as of today date and current timestamp
func (d DynamoDBClient) SetNetworth(username string, networth float64, assets float64, liabilities float64) error {
	now := time.Now().UTC()
	timestamp := now.Format(time.RFC3339)
	networthStr := aws.String(strconv.FormatFloat(networth, 'f', -1, 64))
	assetsStr := aws.String(strconv.FormatFloat(assets, 'f', -1, 64))
	liabilitiesStr := aws.String(strconv.FormatFloat(liabilities, 'f', -1, 64))
	key := fmt.Sprintf("%s:networth", username)

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
							"updated_at":  {S: aws.String(timestamp)},
						},
					},
				},
				{
					PutRequest: &dynamodb.PutRequest{
						Item: map[string]dynamodb.AttributeValue{
							"id":          {S: aws.String(key)},
							"sort":        {S: aws.String("latest")},
							"networth":    {N: networthStr},
							"assets":      {N: assetsStr},
							"liabilities": {N: liabilitiesStr},
							"updated_at":  {S: aws.String(timestamp)},
						},
					},
				},
			},
		},
	})

	_, err := req.Send()

	return err
}

// GetTokens - return all tokens decrypted from db for a username
func (d DynamoDBClient) GetTokens(kms *KMSClient, username string) ([]Token, error) {
	var tokens []Token
	var payload []Token
	key := fmt.Sprintf("%s:token", username)

	req := d.QueryRequest(&dynamodb.QueryInput{
		TableName: aws.String(dbTable),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":id": {S: aws.String(key)},
		},
		KeyConditionExpression: aws.String("id = :id"),
	})

	res, err := req.Send()
	if err != nil {
		return tokens, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &tokens); err != nil {
		return tokens, err
	}

	for _, token := range tokens {
		accessToken, err := kms.Decrypt(token.AccessToken)
		if err != nil {
			log.Printf("Problem decoding access token: %+v\n", err)
			return nil, err
		}
		payload = append(payload, Token{
			AccessToken:   accessToken,
			ItemID:        token.ItemID,
			InstitutionID: token.InstitutionID})
	}

	return payload, nil
}

// GetTokenByItemID - return decrypted token based on item_id
func (d DynamoDBClient) GetTokenByItemID(kms *KMSClient, itemID string) (Token, error) {
	var tokens []Token
	var token Token

	req := d.ScanRequest(&dynamodb.ScanInput{
		TableName:        aws.String(dbTable),
		FilterExpression: aws.String("contains(id, :token) and item_id = :itemID"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":token":  {S: aws.String(":token")},
			":itemID": {S: aws.String(itemID)},
		},
	})

	res, err := req.Send()
	if err != nil {
		return token, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &tokens); err != nil {
		return token, err
	}

	for _, token := range tokens {
		accessToken, err := kms.Decrypt(token.AccessToken)
		if err != nil {
			log.Printf("Problem decoding access token: %+v\n", err)
			return token, err
		}
		token.AccessToken = accessToken
		break
	}

	return token, nil
}

// GetUsernameByItemID - return username based on item_id
func (d DynamoDBClient) GetUsernameByItemID(itemID string) (string, error) {
	var tokens []Token

	req := d.ScanRequest(&dynamodb.ScanInput{
		TableName:        aws.String(dbTable),
		FilterExpression: aws.String("contains(id, :token) and item_id = :itemID and attribute_exists(username)"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":token":  {S: aws.String(":token")},
			":itemID": {S: aws.String(itemID)},
		},
	})

	res, err := req.Send()
	if err != nil {
		return "", err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &tokens); err != nil {
		return "", err
	}

	for _, token := range tokens {
		return token.Username, nil
	}

	return "", nil
}

// SetToken save token to db
func (d DynamoDBClient) SetToken(username string, token *Token) error {
	tokenAttr, err := dynamodbattribute.MarshalMap(token)
	if err != nil {
		log.Printf("Problem marshalling token struct into dyno format: %+v\n", err)
		return err
	}

	dbKey := map[string]dynamodb.AttributeValue{
		"id":       {S: aws.String(fmt.Sprintf("%s:token", username))},
		"sort":     {S: aws.String(token.ItemID)},
		"username": {S: aws.String(username)},
	}

	for k, v := range dbKey {
		tokenAttr[k] = v
	}

	req := d.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(dbTable),
		Item:      tokenAttr,
	})

	if _, err := req.Send(); err != nil {
		log.Printf("Problem SetToken: %+v\n", err)
		return err
	}

	return nil
}

// SetWebhook save webhook to db
func (d DynamoDBClient) SetWebhook(webhook Webhook) error {
	dbAttr, err := dynamodbattribute.MarshalMap(webhook)
	if err != nil {
		log.Printf("Problem marshalling webhook struct into dyno format: %+v\n", err)
		return err
	}

	dbKey := map[string]dynamodb.AttributeValue{
		"id":   {S: aws.String("webhook")},
		"sort": {S: aws.String(webhook.ItemID)},
	}

	for k, v := range dbKey {
		dbAttr[k] = v
	}

	req := d.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(dbTable),
		Item:      dbAttr,
	})

	if _, err := req.Send(); err != nil {
		log.Printf("Problem SetWebhook: %+v\n", err)
		return err
	}

	return nil
}

// SetTransaction save transaction to db
func (d DynamoDBClient) SetTransaction(username string, transaction plaid.Transaction) error {

	transactionAttr, err := dynamodbattribute.MarshalMap(transaction)
	if err != nil {
		log.Printf("Problem marshalling transaction struct into dyno format: %+v\n", err)
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
		log.Printf("Problem SetTransaction: %+v\n", err)
		return err
	}

	return nil
}

// SetAccount save account to db
func (d DynamoDBClient) SetAccount(username string, itemID string, account *plaid.Account) error {
	accountAttr, err := dynamodbattribute.MarshalMap(account)
	if err != nil {
		log.Printf("Problem marshalling account struct into dyno format: %+v\n", err)
		return err
	}

	dbKey := map[string]dynamodb.AttributeValue{
		"id":      {S: aws.String(fmt.Sprintf("%s:account", username))},
		"sort":    {S: aws.String(account.AccountID)},
		"item_id": {S: aws.String(itemID)},
	}

	for k, v := range dbKey {
		accountAttr[k] = v
	}

	req := d.PutItemRequest(&dynamodb.PutItemInput{
		TableName: aws.String(dbTable),
		Item:      accountAttr,
	})

	if _, err := req.Send(); err != nil {
		log.Printf("Problem saving account to db: %+v\n", err)
		return err
	}

	return nil
}

// GetAccounts return all accounts from db for a username
func (d DynamoDBClient) GetAccounts(username string) ([]Account, error) {
	var accounts []Account
	key := fmt.Sprintf("%s:account", username)

	req := d.QueryRequest(&dynamodb.QueryInput{
		TableName: aws.String(dbTable),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":id": {S: aws.String(key)},
		},
		KeyConditionExpression: aws.String("id = :id"),
	})

	res, err := req.Send()
	if err != nil {
		return accounts, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &accounts); err != nil {
		return accounts, err
	}

	return accounts, nil
}

// GetTransactions return all transactions from db for a username
func (d DynamoDBClient) GetTransactions(username string, accountID string) ([]Transaction, error) {
	var transactions []Transaction

	req := d.QueryRequest(&dynamodb.QueryInput{
		TableName:              aws.String(dbTable),
		KeyConditionExpression: aws.String("id = :id and account_id = :accountID"),
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":id":        {S: aws.String(fmt.Sprintf("%s:transaction", username))},
			":accountID": {S: aws.String(accountID)},
		},
	})

	res, err := req.Send()
	if err != nil {
		return transactions, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &transactions); err != nil {
		return transactions, err
	}

	return transactions, nil
}

// GetAllUsers - get all users
func (d DynamoDBClient) GetAllUsers() ([]Token, error) {
	// TODO: Query on username index instead of Scan
	var tokens []Token
	req := d.ScanRequest(&dynamodb.ScanInput{
		TableName:        aws.String(dbTable),
		FilterExpression: aws.String("attribute_exists(username)"),
	})

	res, err := req.Send()
	if err != nil {
		return tokens, err
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(res.Items, &tokens); err != nil {
		return tokens, err
	}

	return tokens, nil
}
