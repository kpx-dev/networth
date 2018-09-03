package main

import (
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleHealthcheck() http.HandlerFunc {
	version := nwlib.GetAPIVersion()

	return func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{"version": version}
		nwlib.SuccessResp(w, payload)
	}
}
