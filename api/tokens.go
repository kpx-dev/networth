package main

import (
	"net/http"
	"strings"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func (s *server) handleTokens() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// switch r.Method {
		// case "GET":
		// case "POST":
		// }

		type claim struct {
			Username string `json:"username"`
		}

		key := []byte(jwtSecret)
		sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT"))
		if err != nil {
			error(w, err.Error())
			return
		}

		myClaim := claim{"demo@networth.app"}

		raw, err := jwt.Signed(sig).Claims(myClaim).CompactSerialize()
		if err != nil {
			error(w, err.Error())
			return
		}

		success(w, raw)
	}
}

func (s *server) auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		jwtKey := strings.Replace(authHeader, "Bearer", "", 1)
		parsed, err := jwt.ParseSigned(jwtKey)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			error(w, "Invalid JWT format")
			return
		}

		claim := jwt.Claims{}
		key := []byte(jwtSecret)
		if err := parsed.Claims(key, &claim); err != nil {
			w.WriteHeader(http.StatusForbidden)
			error(w, "Invalid JWT crypto")
			return
		}

		h(w, r)
	}
}
