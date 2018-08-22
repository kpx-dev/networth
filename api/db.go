package main

// DBClient client for db storage
type DBClient struct {
	*DynamoDBClient
}

// NewDBClient new db client
func NewDBClient() *DBClient {
	client := NewDynamoDBClient()

	return &DBClient{client}
}
