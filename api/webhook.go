package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/networth-app/networth/api/lib"
)

var snsARN = nwlib.GetEnv("SNS_TOPIC_ARN")

func (s *NetworthAPI) handleWebhook() http.HandlerFunc {
	// type WebhookError struct {
	// 	DisplayMessage string `json:"display_message"` //: "The provided credentials were not correct. Please try again.",
	// 	ErrorCode      string `json:"error_code"`      //: "ITEM_LOGIN_REQUIRED",
	// 	ErrorMessage   string `json:"error_message"`   //: "the provided credentials were not correct",
	// 	ErrorType      string `json:"error_type"`      //: "ITEM_ERROR",
	// 	Status         int    `json:"status"`          //: 400
	// }

	// type WebhookBody struct {
	// 	WebhookType         string       `json:"webhook_type"`         // "TRANSACTIONS",
	// 	WebhookCode         string       `json:"webhook_code"`         // INITIAL_UPDATE HISTORICAL_UPDATE DEFAULT_UPDATE TRANSACTIONS_REMOVED
	// 	ItemID              string       `json:"item_id"`              // "wz666MBjYWTp2PDzzggYhM6oWWmBb",
	// 	Error               WebhookError `json:"error"`                // null,
	// 	NewTransactions     int          `json:"new_transactions"`     //19
	// 	RemovedTransactions []string     `json:"removed_transactions"` // ["yBVBEwrPyJs8GvR77N7QTxnGg6wG74H7dEDN6", "kgygNvAVPzSX9KkddNdWHaVGRVex1MHm3k9no"],
	// }

	return func(w http.ResponseWriter, r *http.Request) {
		var webhook nwlib.Webhook
		// Plaid webhook ips: https://support.plaid.com/customer/en/portal/articles/2546264-webhook-overview
		plaidIPs := []string{"52.21.26.131", "52.21.47.157", "52.41.247.19", "52.88.82.239"}
		ips := r.Header.Get("X-Forwarded-For")

		validIP := false
		for ip := range ips {
			for plaidIP := range plaidIPs {
				if ip == plaidIP {
					validIP = true
					break
				}
			}
		}

		if !validIP {
			log.Printf("Invalid webhook, IP does not match whitelisted IP: %+v", ips)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
			nwlib.ErrorResp(w, err.Error())
			return
		}

		// TODO: handle webhook to sync new trans
		nwlib.PublishSNS(snsARN, fmt.Sprintf("New webhook, type: %s, code: %s, item: %s, new trans: %+v", webhook.WebhookType, webhook.WebhookCode, webhook.ItemID, webhook.NewTransactions))
		s.db.SetWebhook(webhook)
		nwlib.SuccessResp(w, webhook)
	}
}
