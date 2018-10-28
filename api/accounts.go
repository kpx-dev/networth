package main

import (
	"log"
	"net/http"

	"github.com/networth-app/networth/lib"
)

// AccountResp - account payload
type AccountResp struct {
	Accounts        []nwlib.Account `json:"accounts"`
	InstitutionName string          `json:"institution_name"`
}

func (s *NetworthAPI) handleAccounts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := make(map[string]AccountResp)

		accounts, err := s.db.GetAccounts(username)

		if err != nil {
			log.Printf("Problem getting accounts: %+v\n", err)
			nwlib.ErrorResp(w, err.Error())
			return
		}

		for _, account := range accounts {
			insID := account.InstitutionID
			insName := account.InstitutionName

			if insID == "" {
				insID = "unknown"
			}

			if insName == "" {
				insName = "Unknown"
			}

			existingAccounts := make([]nwlib.Account, 1)
			if _, ok := payload[insID]; ok {
				existingAccounts = payload[insID].Accounts
			}

			payload[insID] = AccountResp{
				Accounts:        append(existingAccounts, account),
				InstitutionName: insName,
			}
		}

		nwlib.SuccessResp(w, payload)
	}
}
