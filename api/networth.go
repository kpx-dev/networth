package main

import (
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		networth := s.db.GetNetworth(username)

		nwlib.SuccessResp(w, networth)
	}
}
