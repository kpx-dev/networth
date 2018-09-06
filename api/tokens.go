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
		encryptedToken := kmsClient.Encrypt(exchangedToken.AccessToken)

		jwtUsername := s.username(r.Header)
		token := nwlib.Token{
			ItemID:          exchangedToken.ItemID,
			AccessToken:     encryptedToken,
			AccountID:       body.AccountID,
			InstitutionID:   body.InstitutionID,
			InstitutionName: body.InstitutionName,
			Accounts:        body.Accounts,
		}

		tokens := &nwlib.Tokens{
			Tokens: []nwlib.Token{
				token,
			},
		}

		existingTokens := s.db.GetToken(jwtUsername, "")

		for existingInstitutionID := range existingTokens {
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

		if err := s.db.SetToken(jwtUsername, tokenStore); err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		nwlib.SuccessResp(w, "access token created")
	}
}

// func (s *NetworthAPI) handleTokens() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		type claim struct {
// 			Username string `json:"username"`
// 		}

// 		key := []byte(jwtSecret)
// 		sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
// 		if err != nil {
// 			nwlib.ErrorResp(w, err.Error())
// 			return
// 		}

// 		myClaim := claim{username}

// 		raw, err := jwt.Signed(sig).Claims(myClaim).CompactSerialize()
// 		if err != nil {
// 			nwlib.ErrorResp(w, err.Error())
// 			return
// 		}

// 		nwlib.SuccessResp(w, raw)
// 	}
// }

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

		// authHeader := r.Header.Get("Authorization")
		// jwtKey := strings.Replace(authHeader, "Bearer ", "", 1)
		// parsed, err := jwt.ParseSigned(jwtKey)
		// if err != nil {
		// 	w.WriteHeader(http.StatusForbidden)
		// 	nwlib.ErrorResp(w, "Invalid JWT format")
		// 	return
		// }

		// claim := jwt.Claims{}
		// key := []byte(jwtSecret)
		// if err := parsed.Claims(key, &claim); err != nil {
		// 	w.WriteHeader(http.StatusForbidden)
		// 	nwlib.ErrorResp(w, "Invalid JWT crypto")
		// 	return
		// }

		// h(w, r)
	}
}
