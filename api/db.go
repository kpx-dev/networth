package main

import (
	"github.com/networth-app/networth/lib"
)

// DBClient client for db storage
type DBClient struct {
	*nwlib.DynamoDBClient
}

// NewDBClient new db client
func NewDBClient() *DBClient {
	client := nwlib.NewDynamoDBClient()

	return &DBClient{client}
}
