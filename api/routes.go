package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (s *NetworthAPI) routes() {
	s.router.HandleFunc("/tokens", s.auth(s.handleTokens())).Methods("GET", "POST")
	s.router.HandleFunc("/networth", s.auth(s.handleNetworth())).Methods("GET")
	s.router.HandleFunc("/accounts", s.auth(s.handleAccounts()))
	s.router.HandleFunc("/healthcheck", s.handleHealthcheck()).Methods("GET")
}

func (s *NetworthAPI) handleHealthcheck() http.HandlerFunc {
	version := getAPIVersion()

	return func(w http.ResponseWriter, r *http.Request) {

		success(w, "version: "+version)
	}
}

func (s *NetworthAPI) serve(host string) {
	fmt.Println(host)
	srv := &http.Server{
		Handler:      s.router,
		Addr:         host,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting api service on: ", host)
	log.Fatal(srv.ListenAndServe())
}
