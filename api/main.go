package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gorilla/mux"
)

// NetworthAPI nw api struct
type NetworthAPI struct {
	// db     *RedisClient
	db     *BoltClient
	router *mux.Router
	plaid  *PlaidClient
}

var (
	username       = "demo@networth.app"
	accessToken    string
	jwtSecret      string
	plaidEnv       string
	plaidClientID  string
	plaidSecret    string
	plaidPublicKey string
)

// APIGatewayHandler handle incoming apigateway request
func APIGatewayHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	loadDotEnv()

	accessToken = getEnv("PLAID_ACCESS_TOKEN")
	jwtSecret = getEnv("JWT_SECRET")
	plaidClientID = getEnv("PLAID_CLIENT_ID")
	plaidSecret = getEnv("PLAID_SECRET")
	plaidPublicKey = getEnv("PLAID_PUBLIC_KEY")
	plaidEnv = getEnv("PLAID_ENV", "sandbox")
	apiHost := getEnv("API_HOST", ":8000")

	plaidClient := NewPlaidClient()
	// redisClient := NewRedisClient()
	boltClient := NewBoltClient()

	api := &NetworthAPI{
		// db:     redisClient,
		db:     boltClient,
		router: mux.NewRouter(),
		plaid:  plaidClient,
	}
	api.Start(apiHost)
}

func main() {
	lambda.Start(APIGatewayHandler)
}
