package main

import (
	"log"
	"net/http"
	"time"
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
	srv := &http.Server{
		Handler:      s.router,
		Addr:         host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.init()

	log.Println("Starting api service on: ", host)
	log.Fatal(srv.ListenAndServe())
}
