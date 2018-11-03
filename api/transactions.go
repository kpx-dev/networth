package main

import (
	"log"
	"net/http"

	"github.com/networth-app/networth/lib"
)

func (s *NetworthAPI) handleTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		accountID := url.Get("account_id")

		transactions, err := s.db.GetTransactions(username, accountID)

		if err != nil {
			log.Printf("Problem getting transactions: %+v\n", err)
			nwlib.ErrorResp(w, err.Error())
			return
		}

		nwlib.SuccessResp(w, transactions)
	}
}
