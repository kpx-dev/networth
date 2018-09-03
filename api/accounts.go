package main

import (
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: use real access_token
		accessToken := "test"
		accounts, err := s.plaid.GetAccounts(accessToken)

		if err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		// for _, account := range accounts.Accounts {
		// 	fmt.Println(account)
		// }

		// payload, _ := json.Marshal(accounts.Accounts)
		// fmt.Println(string(payload))
		// json.NewEncoder(w).Encode()
		nwlib.SuccessResp(w, accounts.Accounts)
	}
}
