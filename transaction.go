package zinary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

type transactionResponse struct {
	ID                      string  `json:"id"`
	FiatTax                 int64   `json:"fiatTax"`
	FiatFee                 int64   `json:"fiatFee"`
	ExchangeRate            float64 `json:"exchangeRate"`
	CryptocurrencyTax       int64   `json:"cryptocurrencyTax"`
	CryptocurrencyFee       int64   `json:"cryptocurrencyFee"`
	Cryptcurrency           string  `json:"cryptocurrency"`
	CryptoCurrencyAmount    float64 `json:"cryptocurrencyAmount"`
	FiatAmount              int64   `json:"fiatAmount"`
	FiatCurrency            string  `json:"fiatCurrency"`
	Status                  string  `json:"status"`
	Address                 string  `json:"address"`
	BlockchainConfirmations int64   `json:"blockchainConfirmations"`
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
func (c client) CreateTransaction(body Transaction) (transactionResponse, error) {
	var resp transactionResponse

	response, err := c.makeRequest(http.MethodPost, "/create-transaction", &body)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}

type getTransactionResponse struct {
	ID                      string         `json:"id"`
	UniqueID                string         `json:"uniqueId"`
	Status                  string         `json:"status"`
	CryptocurrencyAmount    float64        `json:"cryptocurrencyAmount"`
	Cryptocurrency          string         `json:"cryptocurrency"`
	BlockchainConfirmations int64          `json:"blockchainConfirmations"`
	AmountReceived          float64        `json:"amountReceived"`
	ExchangeRate            float64        `json:"exchangeRate"`
	Type                    string         `json:"type"`
	Details                 map[string]any `json:"details"`
	FiatAmount              int64          `json:"fiatAmount"`
	FiatCurrency            string         `json:"fiatCurrency"`
	BlockchainTransactionID string         `json:"blockchainTransactionId"`
	WebhookURLCalled        bool           `json:"webhookUrlCalled"`
}

// GetTransaction
//
// This endpoint allows you to retrieve or poll the status of a transaction. After creating a transaction using the
// create-transaction endpoint, an ID will be provided as part of the response.
// You can use this ID to track the status of your transaction by making a request to this endpoint.
//
// Transaction Statuses:
//
// - Pending: The user has not made a payment yet.
//
// - Confirming: An incoming transaction has been detected on the blockchain.
//
// - Confirmed: The user has completed the payment, and the blockchain has confirmed it.
//
// - Incomplete: The user made a partial payment.
//
// - Excess: The user paid more than the requested amount.
//
// - Cancelled: The payment window expired without receiving a payment.
//
// - Failed: The transaction failed from the user's end after it was broadcast to the blockchain.
func (c client) GetTransaction(transactionID string) (getTransactionResponse, error) {
	var resp getTransactionResponse

	response, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/transaction/%s", transactionID), nil)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}

type listTransactionsResponse struct {
	LastPage     int                           `json:"lastPage"`
	TotalRecords int                           `json:"totalRecords"`
	CurrentPage  int                           `json:"currentPage"`
	HasMorePages bool                          `json:"hasMorePages"`
	Data         []listTransactionResponseData `json:"data"`
}

type listTransactionResponseData struct {
	ID                      string  `json:"id"`
	Type                    string  `json:"type"`
	CreatedAt               string  `json:"createdAt"`
	Cryptocurrency          string  `json:"cryptocurrency"`
	CryptocurrencyAmount    float64 `json:"cryptocurrencyAmount"`
	FiatAmount              int64   `json:"fiatAmount"`
	FiatCurrency            string  `json:"fiatCurrency"`
	AmountReceived          float64 `json:"amountReceived"`
	BlockchainTransactionID string  `json:"blockchainTransactionId"`
	Status                  string  `json:"status"`
	FiatTax                 int64   `json:"fiatTax"`
	FiatFee                 int64   `json:"fiatFee"`
	CryptocurrencyTax       int64   `json:"cryptocurrencyTax"`
	CryptocurrencyFee       int64   `json:"cryptocurrencyFee"`
}

type ListTransactionsRequest struct {
	Limit *int
	Page  *int
}

// ListTransactionsProd
//
// This endpoint returns a paginated list of all transactions created when your app was in production mode.
// Transactions created in development mode (which are simulated or fake) will not appear in this list.
func (c client) ListTransactionsProd(l ListTransactionsRequest) (listTransactionsResponse, error) {
	var resp listTransactionsResponse
	if l.Limit == nil {
		limit := 50
		l.Limit = &limit
	}

	if l.Page == nil {
		page := 1
		l.Page = &page
	}

	queryParams := url.Values{}
	queryParams.Add("limit", fmt.Sprintf("%d", *l.Limit))
	queryParams.Add("page", fmt.Sprintf("%d", *l.Page))

	response, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/transactions?%s", queryParams.Encode()), nil)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}

// ListTransactionsDev
//
// This endpoint returns a paginated list of all transactions created when your app was in development mode.
// Transactions created in production mode (which are real transactions) will not appear in this list.
func (c client) ListTransactionsDev(l ListTransactionsRequest) (listTransactionsResponse, error) {
	var resp listTransactionsResponse
	if l.Limit == nil {
		limit := 50
		l.Limit = &limit
	}

	if l.Page == nil {
		page := 1
		l.Page = &page
	}

	queryParams := url.Values{}
	queryParams.Add("limit", fmt.Sprintf("%d", *l.Limit))
	queryParams.Add("page", fmt.Sprintf("%d", *l.Page))

	response, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/transactions-dev?%s", queryParams.Encode()), nil)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}
