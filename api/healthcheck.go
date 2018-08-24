package main

import (
	"net/http"
)

func (s *NetworthAPI) handleHealthcheck() http.HandlerFunc {
	version := getAPIVersion()

	return func(w http.ResponseWriter, r *http.Request) {
		payload := map[string]string{"version": version}
		successResp(w, payload)
	}
}
