package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/networth-app/networth/lib"
)

func (s *NetworthAPI) handleGetPublicToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		itemID := url.Get("item_id")

		publicToken, err := s.db.GetToken(kmsClient, username, itemID)

		if err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		nwlib.SuccessResp(w, publicToken)
	}
}

func (s *NetworthAPI) handleTokenExchange() http.HandlerFunc {
	// IncomingToken body from api
	type IncomingToken struct {
		AccessToken     string `json:"access_token"`
		AccountID       string `json:"account_id"`
		InstitutionID   string `json:"institution_id"`
		InstitutionName string `json:"institution_name"`
	}

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

		// TODO: enable for testing purposes only
		// publicToken, _ := s.plaid.CreateSandboxPublicToken("ins_1", []string{"transactions"})
		// body.AccessToken = publicToken.PublicToken

		exchangedToken, err := s.plaid.ExchangePublicToken(body.AccessToken)

		if err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		encryptedToken, err := kmsClient.Encrypt(exchangedToken.AccessToken)

		if err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		token := &nwlib.Token{
			ItemID:          exchangedToken.ItemID,
			AccessToken:     encryptedToken,
			InstitutionID:   body.InstitutionID,
			InstitutionName: body.InstitutionName,
		}

		if err := s.db.SetToken(username, token); err != nil {
			log.Printf("Problem saving token to db: %+v\n", err)
			nwlib.ErrorResp(w, err.Error())
			return
		}

		payload := token
		payload.AccessToken = "*redacted*"
		nwlib.SuccessResp(w, payload)
	}
}
