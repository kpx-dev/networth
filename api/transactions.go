package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/networth-app/networth/lib"
)

func (s *NetworthAPI) handleTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		accountID := vars["accountID"]

		transactions, err := s.db.GetTransactions(username, accountID)

		if err != nil {
			log.Printf("Problem getting transactions for account: %s %+v\n", accountID, err)
			nwlib.ErrorResp(w, err.Error())
			return
		}

		if len(transactions) == 0 {
			nwlib.SuccessResp(w, []string{})
			return
		}

		sort.Slice(transactions, func(i, j int) bool {
			return transactions[i].Date > transactions[j].Date
		})

		nwlib.SuccessResp(w, transactions)
	}
}
