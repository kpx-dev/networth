package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/networth-app/networth/api/lib"
	"gopkg.in/square/go-jose.v2/jwt"
)

// IncomingToken body from api
type IncomingToken struct {
	AccessToken     string   `json:"access_token"`
	Accounts        []string `json:"accounts"`
	AccountID       string   `json:"account_id"`
	InstitutionID   string   `json:"institution_id"`
	InstitutionName string   `json:"institution_name"`
}

func (s *NetworthAPI) handleTokenExchange() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var body IncomingToken

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		// validate body
		if body.InstitutionID == "" || body.AccessToken == "" {
			nwlib.ErrorResp(w, "Missing required body fields.")
			return
		}

		// TODO: remove fixture
		publicToken, err := s.plaid.CreateSandboxPublicToken("ins_1", []string{"transactions"})
		if err != nil {
			log.Println("Problem creating sandbox public token ", err)
			return
		}
		body.AccessToken = publicToken.PublicToken

		exchangedToken, err := s.plaid.ExchangePublicToken(body.AccessToken)

		if err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		kmsClient := nwlib.NewKMSClient()
		encryptedToken, err := kmsClient.Encrypt(exchangedToken.AccessToken)

		if err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		jwtUsername := s.username(r.Header)
		token := &nwlib.Token{
			ItemID:          exchangedToken.ItemID,
			AccessToken:     encryptedToken,
			AccountID:       body.AccountID,
			InstitutionID:   body.InstitutionID,
			InstitutionName: body.InstitutionName,
			Accounts:        body.Accounts,
		}

		if err := s.db.SetToken(jwtUsername, body.InstitutionID, token); err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		payload := token
		payload.AccessToken = "*redacted*"
		nwlib.SuccessResp(w, payload)
	}
}

func (s *NetworthAPI) username(headers http.Header) string {
	type CognitoJWT struct {
		Username string `json:"cognito:username"`
		Email    string `json:"email"`
	}

	authHeader := headers.Get("Authorization")
	jwtKey := strings.Replace(authHeader, "Bearer ", "", 1)
	tok, err := jwt.ParseSigned(jwtKey)
	if err != nil {
		log.Println("Problem parsing jwt ", err)
		return ""
	}

	var claim CognitoJWT
	tok.UnsafeClaimsWithoutVerification(&claim)

	return claim.Username
}
