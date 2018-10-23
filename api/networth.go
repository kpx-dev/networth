package main

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/networth-app/networth/lib"
)

func (s *NetworthAPI) handleNetworth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		networth, err := s.db.GetNetworth(username)
		if err != nil {
			nwlib.ErrorResp(w, err.Error())
		}
		nwlib.SuccessResp(w, networth)
	}
}

func (s *NetworthAPI) handleNetworthHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Query()
		startDate := url.Get("start_date")
		endDate := url.Get("end_date")
		resolution := url.Get("resolution")

		// set the start end date if missing
		if resolution != "" && (startDate == "" || endDate == "") {
			now := time.Now()
			endDate = time.Now().Format("2006-01-02")

			switch resolution {
			case "daily":
				// last 1 day:
				startDate = now.AddDate(0, 0, -1).Format("2006-01-02")
				break
			case "weekly", "monthly":
				// last 30 days (1 month)
				startDate = now.AddDate(0, -1, 0).Format("2006-01-02")
				break
			case "yearly":
				// last 12 months
				startDate = now.AddDate(-1, 0, 0).Format("2006-01-02")
				break
			default:
				// last 30 days (1 month)
				startDate = now.AddDate(0, -1, 0).Format("2006-01-02")
			}
		}

		networth, err := s.db.GetNetworthByDateRange(username, startDate, endDate)
		if err != nil {
			nwlib.ErrorResp(w, err.Error())
		}

		grouped := groupBy(networth, resolution)
		nwlib.SuccessResp(w, grouped)
	}
}

func groupBy(networth []nwlib.Networth, resolution string) map[string]float64 {
	payload := make(map[string]float64)
	cache := make(map[string][]float64)

	// map - group net worth value by date resolution
	for _, val := range networth {
		var key string
		ts, _ := time.Parse(time.RFC3339, val.DateTime)
		year, month, day := ts.Date()
		hour := ts.Hour()

		switch resolution {
		case "daily":
			// group by hour:
			key = fmt.Sprintf("%d-%02d-%02dT%02d:00:00Z", year, month, day, hour)
			break
		case "weekly", "monthly":
			// group by day:
			key = fmt.Sprintf("%d-%02d-%02dT00:00:00Z", year, month, day)
			break
		case "yearly":
			// group by month:
			key = fmt.Sprintf("%d-%02d-01T00:00:00Z", year, month)
			break
		default:
			// use weekly, group by day:
			key = fmt.Sprintf("%d-%02d-%02dT00:00:00Z", year, month, day)
		}

		cache[key] = append(cache[key], val.Networth)
	}

	// reduce - average out the net worth value per key
	for k, v := range cache {
		total := 0.0
		for _, eachVal := range v {
			total += eachVal
		}
		average := total / float64(len(v))
		payload[k] = math.Round(average*100) / 100
	}

	return payload
}
