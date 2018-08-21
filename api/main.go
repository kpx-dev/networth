package main

import (
	"encoding/json"
	"log"
	"net/http"

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
// func APIGatewayHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
// 	loadDotEnv()

// 	accessToken = getEnv("PLAID_ACCESS_TOKEN")
// 	jwtSecret = getEnv("JWT_SECRET")
// 	plaidClientID = getEnv("PLAID_CLIENT_ID")
// 	plaidSecret = getEnv("PLAID_SECRET")
// 	plaidPublicKey = getEnv("PLAID_PUBLIC_KEY")
// 	plaidEnv = getEnv("PLAID_ENV", "sandbox")
// 	apiHost := getEnv("API_HOST", ":8000")

// 	plaidClient := NewPlaidClient()
// 	// redisClient := NewRedisClient()
// 	boltClient := NewBoltClient()

// 	api := &NetworthAPI{
// 		// db:     redisClient,
// 		db:     boltClient,
// 		router: mux.NewRouter(),
// 		plaid:  plaidClient,
// 	}
// 	api.Start(apiHost)
// }

func getNetworth(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// The APIGatewayProxyResponse.Body field needs to be a string, so
	// we marshal the book record into JSON.
	js, err := json.Marshal(bk)
	if err != nil {
		return serverError(err)
	}

	// Return a response with a 200 OK status and the JSON book record
	// as the body.
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(js),
	}, nil
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.Path {
	case "networth":
		return getNetworth(req)
	default:
		return clientError(http.StatusBadRequest)
	}
}

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	log.Printf(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
