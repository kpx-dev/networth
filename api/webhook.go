package main

import (
	"encoding/json"
	"net/http"
)

func (s *NetworthAPI) handleWebhook() http.HandlerFunc {
	type WebhookError struct {
		DisplayMessage string `json:"display_message"` //: "The provided credentials were not correct. Please try again.",
		ErrorCode      string `json:"error_code"`      //: "ITEM_LOGIN_REQUIRED",
		ErrorMessage   string `json:"error_message"`   //: "the provided credentials were not correct",
		ErrorType      string `json:"error_type"`      //: "ITEM_ERROR",
		Status         int    `json:"status"`          //: 400
	}

	type WebhookBody struct {
		WebhookType         string       `json:"webhook_type"`         // "TRANSACTIONS",// INITIAL_UPDATE HISTORICAL_UPDATE DEFAULT_UPDATE TRANSACTIONS_REMOVED
		WebhookCode         string       `json:"webhook_code"`         // "INITIAL_UPDATE",
		ItemID              string       `json:"item_id"`              // "wz666MBjYWTp2PDzzggYhM6oWWmBb",
		Error               WebhookError `json:"error"`                // null,
		NewTransactions     int          `json:"new_transactions"`     //19
		RemovedTransactions []string     `json:"removed_transactions"` // ["yBVBEwrPyJs8GvR77N7QTxnGg6wG74H7dEDN6", "kgygNvAVPzSX9KkddNdWHaVGRVex1MHm3k9no"],
	}

	// plaid webhook ips: (https://support.plaid.com/customer/en/portal/articles/2546264-webhook-overview)

	return func(w http.ResponseWriter, r *http.Request) {
		var body WebhookBody
		// plaidIps := []string{
		// 	"52.21.26.131", "52.21.47.157", "52.41.247.19", "52.88.82.239",
		// }

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errorResp(w, err.Error())
			return
		}

		// TODO: check to make sure ip came from whitelist
		// ips := r.Header.Get("X-Forwarded-For")
		// fmt.Println("Got webhook message from these ips: ", ips, r.RemoteAddr)
		alert("New webhook, type:" + body.WebhookType + " code: " + body.WebhookCode)
		successResp(w, body)
	}
}
