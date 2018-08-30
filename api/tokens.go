package main

import (
	"encoding/json"
	"net/http"
	"strings"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func (s *NetworthAPI) handleTokenExchange() http.HandlerFunc {
	type TokenBody struct {
		Token string `json:"token"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var body TokenBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errorResp(w, err.Error())
			return
		}

		token, err := s.plaid.ExchangePublicToken(body.Token)

		if err != nil {
			errorResp(w, err.Error())
			return
		}

		jwtUsername := s.username(r.Header)
		tokenMap := map[string]string{
			"token": token.AccessToken,
		}

		if err := s.db.Set(tokenTable, jwtUsername, "", tokenMap); err != nil {
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
		Username string `json:"username"`
	}

	authHeader := headers.Get("Authorization")
	jwtKey := strings.Replace(authHeader, "Bearer ", "", 1)
	tok, err := jwt.ParseSigned(jwtKey)
	if err != nil {
		return ""
	}

	claim := &CognitoJWT{}
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
