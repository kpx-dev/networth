package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Token hold the structure for saving to db
type Token struct {
	AccessTokens    []string `json:"access_tokens"`
	Accounts        []string `json:"accounts"`
	AccountID       string   `json:"account_id"`
	InstitutionID   string   `json:"institution_id"`
	InstitutionName string   `json:"institution_name"`
}

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
			errorResp(w, err.Error())
			return
		}

		// validate body
		if body.InstitutionID == "" || body.AccessToken == "" {
			errorResp(w, "Missing required body fields.")
			return
		}

		// TODO: remove fixture
		// publicToken, err := s.plaid.CreateSandboxPublicToken("ins_1", []string{"transactions"})
		// body.AccessToken = publicToken.PublicToken

		token, err := s.plaid.ExchangePublicToken(body.AccessToken)

		if err != nil {
			errorResp(w, err.Error())
			return
		}

		kmsClient := NewKMSClient()
		encryptedToken := kmsClient.Encrypt(token.AccessToken)

		jwtUsername := s.username(r.Header)
		tokenStore := &Token{
			InstitutionName: body.InstitutionName,
			AccessTokens:    []string{encryptedToken},
			Accounts:        body.Accounts,
			AccountID:       body.AccountID,
		}

		tokens := s.db.GetToken(jwtUsername)

		for existingInstitutionID := range tokens {
			if existingInstitutionID == body.InstitutionID {
				// TODO: use Token struct intead of interface
				existingToken := tokens[existingInstitutionID]
				tokenMap := existingToken.(map[string]interface{})
				accessTokens := tokenMap["access_tokens"]
				tokensArray := accessTokens.([]interface{})

				for oldToken := range tokensArray {
					tokenStore.AccessTokens = append(tokenStore.AccessTokens, string(tokensArray[oldToken].(string)))
				}
			}
		}

		if err := s.db.SetToken(jwtUsername, body.InstitutionID, tokenStore); err != nil {
			errorResp(w, err.Error())
			return
		}

		successResp(w, "access token created")
	}
}

func (s *NetworthAPI) handleTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type claim struct {
			Username string `json:"username"`
		}

		key := []byte(jwtSecret)
		sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
		if err != nil {
			errorResp(w, err.Error())
			return
		}

		myClaim := claim{username}

		raw, err := jwt.Signed(sig).Claims(myClaim).CompactSerialize()
		if err != nil {
			errorResp(w, err.Error())
			return
		}

		successResp(w, raw)
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

func (s *NetworthAPI) auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: enable auth when deploy
		h(w, r)
		return

		authHeader := r.Header.Get("Authorization")
		jwtKey := strings.Replace(authHeader, "Bearer ", "", 1)
		parsed, err := jwt.ParseSigned(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			errorResp(w, "Invalid JWT format")
			return
		}

		claim := jwt.Claims{}
		key := []byte(jwtSecret)
		if err := parsed.Claims(key, &claim); err != nil {
			w.WriteHeader(http.StatusForbidden)
			errorResp(w, "Invalid JWT crypto")
			return
		}

		h(w, r)
	}
}
