package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/apex/gateway"
	"github.com/networth-app/networth/api/lib"
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
	s.router.HandleFunc(fmt.Sprintf("%s/networth", prefix), s.handleNetworth()).Methods("GET")
	s.router.HandleFunc(fmt.Sprintf("%s/accounts", prefix), s.handleAccounts())

	s.router.Use(loggingMiddleware)
	s.router.Use(extractUsernameMiddleware)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.SplitAfter(r.URL.String(), prefix)
		log.Println(fmt.Sprintf("%s %s", r.Method, url[1]))
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
			log.Println("Problem parsing jwt ", err)
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

	log.Println(fmt.Sprintf("API service started on: %s", host))

	if nwlib.GetEnv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		log.Fatal(http.ListenAndServe(host, handler))
	} else {
		log.Fatal(gateway.ListenAndServe(host, handler))
	}
}
