package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	type NetworthBody struct {
		Networth string `json:"networth"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			networth := s.db.GetNetworth()
			nwlib.SuccessResp(w, networth)
			break
		case "POST", "PUT":
			var body NetworthBody

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				nwlib.ErrorResp(w, err.Error())
				return
			}

			networth, err := strconv.ParseFloat(body.Networth, 64)

			if err != nil {
				nwlib.ErrorResp(w, err.Error())
				return
			}

			err = s.db.SetNetworth(networth)

			if err != nil {
				nwlib.ErrorResp(w, err.Error())
				return
			}

			nwlib.SuccessResp(w, networth)
			break
		}
	}
}
