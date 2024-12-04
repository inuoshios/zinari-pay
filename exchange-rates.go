package zinary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type exchangeRateResponse struct {
	Rate float64 `json:"rate"`
}

type ExchangeRateRequest struct {
	From string
	To   string
}

func (c client) GetExchangeRate(e ExchangeRateRequest) (exchangeRateResponse, error) {
	var resp exchangeRateResponse

	queryParams := url.Values{}
	queryParams.Add("from", e.From)
	queryParams.Add("to", e.To)

	response, err := c.makeRequest(http.MethodGet, fmt.Sprintf("/exchange-rate?%s", queryParams.Encode()), nil)
	if err != nil {
		_ = json.Unmarshal(response, &resp)
		return resp, err
	}
	_ = json.Unmarshal(response, &resp)
	return resp, nil
}
