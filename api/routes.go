package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/networth-app/networth/api/lib"
	"github.com/rs/cors"
)

func (s *NetworthAPI) init() {
	prefix := "/api"

	s.router.HandleFunc(fmt.Sprintf("%s/tokens/exchange", prefix), s.auth(s.handleTokenExchange())).Methods("POST")
	s.router.HandleFunc(fmt.Sprintf("%s/networth", prefix), s.auth(s.handleNetworth())).Methods("GET", "POST", "PUT")
	s.router.HandleFunc(fmt.Sprintf("%s/webhook", prefix), s.auth(s.handleWebhook())).Methods("POST")
	s.router.HandleFunc(fmt.Sprintf("%s/accounts", prefix), s.auth(s.handleAccounts()))
	s.router.HandleFunc(fmt.Sprintf("%s/healthcheck", prefix), s.handleHealthcheck()).Methods("GET")
	s.router.HandleFunc(prefix, s.handleHealthcheck()).Methods("GET")
	s.router.HandleFunc("/", s.handleHealthcheck()).Methods("GET")
}

// Start start api service
func (s *NetworthAPI) Start(host string) {
	s.init()
	handler := cors.Default().Handler(s.router)

	log.Println("Starting api service on: ", host)

	if nwlib.GetEnv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		log.Fatal(http.ListenAndServe(host, handler))
	} else {
		log.Fatal(gateway.ListenAndServe(host, handler))
	}
}
