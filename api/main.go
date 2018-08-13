package main

import (
	"github.com/gorilla/mux"
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
