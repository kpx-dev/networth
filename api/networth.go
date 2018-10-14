package main

import (
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		startDate := url.Get("start_date")
		endDate := url.Get("end_date")

		if startDate != "" && endDate != "" {
			networth, err := s.db.GetNetworthByDateRange(username, startDate, endDate)
			if err != nil {
				nwlib.ErrorResp(w, err.Error())
			}
			nwlib.SuccessResp(w, networth)

		} else {
			networth, err := s.db.GetNetworth(username)
			if err != nil {
				nwlib.ErrorResp(w, err.Error())
			}
			nwlib.SuccessResp(w, networth)
		}
	}
}
