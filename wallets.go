package zinary

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WithdrawFromWalletRequest struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

type withdrawFromWalletResponse struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type getWalletsResponse struct {
	ID       string  `json:"id"`
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

// GetWallets
//
// This endpoint allows you to retrieve a list of all your app wallets,
// including their balances and wallet IDs. The wallet ID provided can be used for withdrawals.
func (c client) GetWallets() []getWalletsResponse {
	var resp []getWalletsResponse

	response, err := c.makeRequest(http.MethodGet, "/wallets", nil)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp
	}
	_ = json.Unmarshal(response, &resp)
	return append([]getWalletsResponse{}, resp...)
}

// GetWalletByID
//
// This endpoint allows you to retrieve the details of a specific wallet by its ID.
// It returns information such as the wallet's balance and any other associated details.
// The wallet ID is typically obtained from the Get Wallets endpoint and can be used to
// track a specific wallet or perform actions like withdrawals.
func (c client) GetWalletByID(walletID string) (getWalletsResponse, error) {
	var resp getWalletsResponse

	response, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/wallet/%s", walletID), nil)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}

// WithdrawFromWallet
//
// This endpoint allows you to initiate a withdrawal from a specific wallet.
// You need to provide the wallet ID (which can be retrieved from the Get Wallets endpoint)
// along with the withdrawal amount and destination address. Once the request is made,
// the system processes the withdrawal, and you can track its status through the appropriate endpoints.
func (c client) WithdrawFromWallet(walletID string, body WithdrawFromWalletRequest) (withdrawFromWalletResponse, error) {
	var resp withdrawFromWalletResponse

	response, err := c.makeRequest(http.MethodPost, fmt.Sprintf("/withdraw/%s", walletID), &body)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}
