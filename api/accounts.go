package main

import (
	"log"
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.db.GetAccounts(username)

		if err != nil {
			log.Println("Problem getting accounts ", err)
			nwlib.ErrorResp(w, err.Error())
			return
		}

		nwlib.SuccessResp(w, accounts)
	}
}
