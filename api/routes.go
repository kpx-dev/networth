package main

import (
	"log"
	"net/http"
)

// init routing
func (s *NetworthAPI) init() {
	s.router.HandleFunc("/tokens/exchange", s.auth(s.handleTokenExchange())).Methods("POST")
	s.router.HandleFunc("/tokens", s.auth(s.handleTokens())).Methods("GET", "POST")
	s.router.HandleFunc("/networth", s.auth(s.handleNetworth())).Methods("GET")
	s.router.HandleFunc("/accounts", s.auth(s.handleAccounts()))
	s.router.HandleFunc("/healthcheck", s.handleHealthcheck()).Methods("GET")
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("../ui/")))
}

// Start start api service
func (s *NetworthAPI) Start(host string) {
	s.init()
	log.Println("Starting api service on: ", host)
	// log.Fatal(gateway.ListenAndServe(host, s.router))
	log.Fatal(http.ListenAndServe(host, s.router))
}
