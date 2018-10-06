package main

import (
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		networth, err := s.db.GetNetworth(username)

		if err != nil {
			nwlib.ErrorResp(w, err.Error())
		}

		nwlib.SuccessResp(w, networth)
	}
}
