package main

import (
	"net/http"
)

func (s *NetworthAPI) handleAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := s.plaid.GetAccounts(accessToken)

		if err != nil {
			error(w, err.Error())
			return
		}

		// for _, account := range accounts.Accounts {
		// 	fmt.Println(account)
		// }

		// payload, _ := json.Marshal(accounts.Accounts)
		// fmt.Println(string(payload))
		// json.NewEncoder(w).Encode()
		success(w, accounts.Accounts)
	}
}
