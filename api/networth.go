package main

import (
	"net/http"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		networth := s.db.GetNetworth()

		successResp(w, networth)
	}
}
