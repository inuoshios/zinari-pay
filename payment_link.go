package zinary

import (
	"encoding/json"
	"net/http"
)

type PaymentLink struct {
	// The fiat currency amount set in your app.
	FiatAmount int64 `json:"fiatAmount"`

	// The email address of the customer making the payment. You will receive notifications for your own records.
	NotificationEmailAddress string `json:"notificationEmailAddress"`

	// An object containing unique identifiers for your payment, such as userId or productId. This object will be included in the webhook notification payload with any payment updates
	Details map[string]any `json:"details"`

	// A URI for redirecting the customer upon successful payment, which can also be set in your ZinariPay console.
	SuccessRedirectURI string `json:"successRedirectUri"`

	// A URI for redirecting the customer upon failed payment, which can also be set in your ZinariPay console.
	FailureRedirectURI string `json:"failureRedirectUri"`
}

// GetPaymentLink
//
// This endpoint allows developers to generate a payment link, which they can redirect customers
// to for completing a purchase. With this Paylink, you can offer your customers a complete payment
// experience without needing to build any UI elements yourself.
// ZinariPay will guide your customers through the payment process and return them to your app using
// the successRedirectUri or failureRedirectUri provided
func (c client) GetPaymentLink(body PaymentLink) (map[string]any, error) {
	var resp map[string]any

	response, err := c.makeRequest(http.MethodPost, "/paylink", &body)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}
