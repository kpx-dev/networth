package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	type NetworthBody struct {
		Networth string `json:"networth"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			networth := s.db.GetNetworth()
			successResp(w, networth)
			break
		case "POST", "PUT":
			var body NetworthBody

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				errorResp(w, err.Error())
				return
			}

			networth, err := strconv.ParseFloat(body.Networth, 64)

			if err != nil {
				errorResp(w, err.Error())
				return
			}

			err = s.db.SetNetworth(networth)

			if err != nil {
				errorResp(w, err.Error())
				return
			}

			successResp(w, networth)
			break
		}
	}
}
