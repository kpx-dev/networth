package main

import (
	"encoding/json"
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
	s.router.HandleFunc("/institutions", s.auth(s.handleInstitutions()))
	s.router.HandleFunc("/healthcheck", s.handleHealthcheck())
	s.router.HandleFunc("/", s.handleIndex())
}

func (s *server) auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// if !currentUser(r).IsAdmin {
		// 	http.NotFound(w, r)
		// 	return
		// }
		h(w, r)
	}
}

func (s *server) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := resp{"root"}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
	}
}

func (s *server) handleInstitutions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		ins, _ := s.plaid.GetInstitutions(5, 0)
		// fmt.Println(ins)

		json.NewEncoder(w).Encode(ins)
	}
}

func (s *server) handleHealthcheck() http.HandlerFunc {
	version := getAPIVersion()

	return func(w http.ResponseWriter, r *http.Request) {
		message := map[string]string{
			"message": "ok!",
			"version": version,
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(message)
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
