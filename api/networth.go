package main

import (
	"net/http"
	"strconv"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		networthStr := s.db.GetNetworth()

		networth, _ := strconv.ParseFloat(networthStr, 32)
		success(w, networth)
	}
}
