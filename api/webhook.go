package main

import (
	"encoding/json"
	"fmt"
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

	return func(w http.ResponseWriter, r *http.Request) {
		var body WebhookBody

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errorResp(w, err.Error())
			return
		}

		fmt.Println("webhook body is")
		fmt.Println(body.WebhookType, body.WebhookCode, body.ItemID)

		successResp(w, body)
	}
}
