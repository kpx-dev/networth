package main

import (
	"github.com/gorilla/mux"
)

var (
	jwtSecret      = getEnv("JWT_SECRET", "FIRE!")
	plaidEnv       = getEnv("PLAID_ENV", "sandbox")
	plaidClientID  = getEnv("PLAID_CLIENT_ID")
	plaidSecret    = getEnv("PLAID_SECRET")
	plaidPublicKey = getEnv("PLAID_PUBLIC_KEY")
)

func main() {
	loadDotEnv()

	plaidClient := NewPlaidClient()
	redisClient := NewRedisClient()

	s := &server{
		db:     redisClient,
		router: mux.NewRouter(),
		plaid:  plaidClient,
	}
	s.routes()
	s.serve(":8000")
}
