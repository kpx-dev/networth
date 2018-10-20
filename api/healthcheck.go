package main

import (
	"net/http"

	"github.com/networth-app/networth/lib"
)

func (s *NetworthAPI) handleHealthcheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{"status": "ok"}
		nwlib.SuccessResp(w, payload)
	}
}
