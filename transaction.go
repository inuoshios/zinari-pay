package zinary

import (
	"encoding/json"
	"net/http"
)

type Transaction struct {
	// Only USDT and USDC are allowed for the cryptocurrency field.
	CryptoCurrency string `json:"cryptocurrency"`
	// The fiatAmount must be greater than $5 in value and represents the fiat amount set in your app, not the cryptocurrency
	// value. This fiat amount will be automatically converted to the appropriate cryptocurrency value. In this example the fiat
	// currency is Nigerian Naira.
	FiatAmount int64 `json:"fiatAmount"`
	// The notificationEmailAddress is the email address of the customer, who will receive notifications about the payment. As the
	// developer or business owner, you are already set up to receive your own notifications.
	NotificationEmailAddress string `json:"notificationEmailAddress"`
	// Any additional properties you need passed in the request should be included in the details object, which will be sent along
	// with the webhook notification when the transaction is processed. This allows you to pass identifiers such as userId,
	// orderId, productId, or any other information needed to track the transaction on your end.
	Details map[string]any `json:"details"`
}

// CreateTransaction
//
// This endpoint allows you to generate a payment address for your user. In development mode, transactions will automatically
// transition to the confirmed state after a few seconds, and webhook URLs will be triggered with an extra isFake: true attribute.
//
// Key Details:
//
// Timeout and Cancellation: If no payment is made within 10 minutes, the transaction will automatically transition to Cancelled.
// After this, payments made to the supplied address will result in lost funds.
// The required arguments for this endpoint are cryptocurrency, fiatAmount, and notificationEmailAddress.
func (c client) CreateTransaction(body Transaction) (map[string]interface{}, error) {
	var resp map[string]any

	response, err := c.makeRequest(http.MethodPost, "/create-transaction", &body)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}
