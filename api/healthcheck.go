package main

import (
	"net/http"
)

func (s *NetworthAPI) handleHealthcheck() http.HandlerFunc {
	version := getAPIVersion()

	return func(w http.ResponseWriter, r *http.Request) {
		successResp(w, "version: "+version)
	}
}
