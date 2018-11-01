package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/apex/gateway"
	"github.com/networth-app/networth/lib"
	"github.com/rs/cors"
	"gopkg.in/square/go-jose.v2/jwt"
)

var prefix = "/api"

func (s *NetworthAPI) init() {
	// unauth routes
	s.router.HandleFunc(fmt.Sprintf("%s/tokens", prefix), s.handleTokenExchange()).Methods("POST")
	s.router.HandleFunc(fmt.Sprintf("%s/webhook", prefix), s.handleWebhook()).Methods("POST")
	s.router.HandleFunc(fmt.Sprintf("%s/healthcheck", prefix), s.handleHealthcheck()).Methods("GET")

	// auth routes
	s.router.HandleFunc(fmt.Sprintf("%s/tokens/public", prefix), s.handleGetPublicToken()).Methods("GET")
	s.router.HandleFunc(fmt.Sprintf("%s/networth", prefix), s.handleNetworth()).Methods("GET")
	s.router.HandleFunc(fmt.Sprintf("%s/networth_history", prefix), s.handleNetworthHistory()).Methods("GET")
	s.router.HandleFunc(fmt.Sprintf("%s/accounts", prefix), s.handleAccounts()).Methods("GET")
	s.router.HandleFunc(fmt.Sprintf("%s/transactions", prefix), s.handleTransactions()).Methods("GET")
	s.router.HandleFunc(fmt.Sprintf("%s/ws", prefix), s.handleWebSocket()).Methods("GET")

	s.router.Use(loggingMiddleware)
	s.router.Use(extractUsernameMiddleware)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.SplitAfter(r.URL.String(), prefix)
		log.Printf("%s %s\n", r.Method, url[1])
		next.ServeHTTP(w, r)
	})
}

func extractUsernameMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type CognitoJWT struct {
			Username string `json:"cognito:username"`
			Email    string `json:"email"`
		}

		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			log.Println("No Authorization header found")
			next.ServeHTTP(w, r)
			return
		}

		jwtKey := strings.Replace(authHeader, "Bearer ", "", 1)

		tok, err := jwt.ParseSigned(jwtKey)
		if err != nil {
			log.Printf("Problem parsing jwt: %+v\n", err)
			next.ServeHTTP(w, r)
			return
		}

		var claim CognitoJWT
		tok.UnsafeClaimsWithoutVerification(&claim)

		username = claim.Username

		next.ServeHTTP(w, r)
	})
}

// Start start api service
func (s *NetworthAPI) Start(host string) {
	s.init()
	handler := cors.Default().Handler(s.router)
	log.Printf("REST API and Websocket service started on: %s\n", host)

	if nwlib.GetEnv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		log.Fatal(http.ListenAndServe(host, handler))
	} else {
		log.Fatal(gateway.ListenAndServe(host, handler))
	}
}
