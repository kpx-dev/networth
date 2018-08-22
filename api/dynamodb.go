package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	tokenTable       = aws.String(getEnv("TOKEN_TABLE"))
	transactionTable = aws.String(getEnv("TRANSACTION_TABLE"))
)

// DynamoDBClient db client struct
type DynamoDBClient struct {
	*dynamodb.DynamoDB
}

// NewDynamoDBClient new dynamodb client
func NewDynamoDBClient() *DynamoDBClient {
	cfg := loadAWSConfig()
	table := dynamodb.New(cfg)

	return &DynamoDBClient{table}
}

// GetNetworth return networth
func (d DynamoDBClient) GetNetworth() float64 {
	// lastHour := time.Now().UTC().Add(-1 * time.Hour).Format(time.RFC3339)
	lastHour := time.Now().UTC().Add(-48 * time.Hour).Format(time.RFC3339)
	now := time.Now().UTC().Format(time.RFC3339)

	req := d.QueryRequest(&dynamodb.QueryInput{
		KeyConditionExpression: aws.String("email=:email AND #sort BETWEEN :lastHour AND :now"),
		ExpressionAttributeNames: map[string]string{
			"#sort": "datetime",
		},
		ExpressionAttributeValues: map[string]dynamodb.AttributeValue{
			":email":    {S: aws.String(username)},
			":lastHour": {S: aws.String(lastHour)},
			":now":      {S: aws.String(now)},
		},
		Limit:     aws.Int64(1),
		TableName: aws.String(getEnv("HISTORY_TABLE")),
	})

	res, err := req.Send()
	if err != nil {
		panic(err)
	}

	if *res.Count > int64(0) {
		// nw := make(map[string]interface{})
		// if err := dynamodbattribute.UnmarshalMap(res.Items[0], &nw); err != nil {
		// 	panic(err)
		// }

		// fmt.Println(nw)
		fmt.Println(res.Items[0]["networth"])
		return 1
	}

	fmt.Println(res)
	return 0.0
}

// Set key / val to db
func (d DynamoDBClient) Set(key string, value string) error {
	fmt.Println("saving ", key, value)
	return nil
}

// GetTokens return tokens from db
// func (d DynamoDBClient) GetTokens(username string) []string {
// 	req := d.GetItemRequest(&dynamodb.GetItemInput{
// 		TableName: accountTable,
// 		Key: map[string]dynamodb.AttributeValue{
// 			"username": {S: aws.String(fmt.Sprintf("%s:tokens", username))},
// 		},
// 	})

// 	res, err := req.Send()
// 	if err != nil {
// 		panic(err.Error())
// 	}

// 	tokens := make(map[string]interface{})
// 	if err := dynamodbattribute.UnmarshalMap(res.Item, &tokens); err != nil {
// 		panic(err)
// 	}

// 	// TODO: might not be a right approach to init kms client here...
// 	kmsClient := NewKMSClient()

// 	payload := []string{""}
// 	for k, v := range tokens {
// 		if strings.HasPrefix(k, "ins_") {
// 			decrypted := kmsClient.Decrypt(v.([]string)[0])
// 			payload = append(payload, decrypted)
// 		}
// 	}

// 	return payload
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
