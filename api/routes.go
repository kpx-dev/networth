package main

import (
	"log"

	"github.com/apex/gateway"
	"github.com/rs/cors"
)

// init routing
func (s *NetworthAPI) init() {
	s.router.HandleFunc("/tokens/exchange", s.auth(s.handleTokenExchange())).Methods("POST")
	s.router.HandleFunc("/tokens", s.auth(s.handleTokens())).Methods("GET", "POST")
	s.router.HandleFunc("/networth", s.auth(s.handleNetworth())).Methods("GET", "POST", "PUT")
	s.router.HandleFunc("/webhook", s.auth(s.handleWebhook())).Methods("POST")
	s.router.HandleFunc("/accounts", s.auth(s.handleAccounts()))
	s.router.HandleFunc("/healthcheck", s.handleHealthcheck()).Methods("GET")
	s.router.HandleFunc("/", s.handleHealthcheck()).Methods("GET")
}

// Start start api service
func (s *NetworthAPI) Start(host string) {
	s.init()
	handler := cors.Default().Handler(s.router)

	log.Println("Starting api service on: ", host)

	log.Fatal(gateway.ListenAndServe(host, handler))
	// log.Fatal(http.ListenAndServe(host, handler))
}
