package main

import (
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
	accessToken    string
	jwtSecret      string
	plaidEnv       string
	plaidClientID  string
	plaidSecret    string
	plaidPublicKey string
)

func main() {
	loadDotEnv()

	accessToken = getEnv("PLAID_ACCESS_TOKEN")
	jwtSecret = getEnv("JWT_SECRET")
	plaidClientID = getEnv("PLAID_CLIENT_ID")
	plaidSecret = getEnv("PLAID_SECRET")
	plaidPublicKey = getEnv("PLAID_PUBLIC_KEY")
	plaidEnv = getEnv("PLAID_ENV", "sandbox")

	plaidClient := NewPlaidClient()
	// redisClient := NewRedisClient()
	boltClient := NewBoltClient()

	s := &NetworthAPI{
		// db:     redisClient,
		db:     boltClient,
		router: mux.NewRouter(),
		plaid:  plaidClient,
	}
	s.routes()
	s.serve(":8000")
}
