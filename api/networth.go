package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtUsername := s.username(r.Header)

		switch r.Method {
		case "GET":
			fmt.Println("username is ", jwtUsername)

			networth := s.db.GetNetworth(jwtUsername)
			nwlib.SuccessResp(w, networth)
			break
		case "POST", "PUT":
			var body nwlib.Networth

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				nwlib.ErrorResp(w, err.Error())
				return
			}

			err := s.db.SetNetworth(jwtUsername, body.Networth, 0.0, 0.0)

			if err != nil {
				nwlib.ErrorResp(w, err.Error())
				return
			}

			nwlib.SuccessResp(w, body.Networth)
			break
		}
	}
}
