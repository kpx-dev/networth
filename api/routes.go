package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	db     *RedisClient
	router *mux.Router
	plaid  *PlaidClient
}

type resp struct {
	Message string `json:"message"`
}

func (s *server) routes() {
	// s.router.HandleFunc("/tokens", s.handleTokens()).Methods("GET", "POST")
	s.router.HandleFunc("/tokens", s.auth(s.handleTokens())).Methods("GET", "POST")
	s.router.HandleFunc("/institutions", s.auth(s.handleInstitutions()))
	s.router.HandleFunc("/healthcheck", s.handleHealthcheck()).Methods("GET")
	s.router.HandleFunc("/", s.handleIndex()).Methods("GET")
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := resp{"networth-api! Please visit https://docs.networth.app/"}

		success(w, message)
	}
}

func (s *server) handleInstitutions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ins, _ := s.plaid.GetInstitutions(5, 0)
		// success(w, ins)
		success(w, "ok")
	}
}

func (s *server) handleHealthcheck() http.HandlerFunc {
	version := getAPIVersion()

	return func(w http.ResponseWriter, r *http.Request) {
		message := map[string]string{
			"message": "ok!",
			"version": version,
		}

		success(w, message)
	}
}

func (s *server) serve(host string) {
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
