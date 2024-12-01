package zinary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	BaseURL = ""

	// DefaultHTTPTimeout is the default timeout on the http client
	DefaultHTTPTimeout = 60 * time.Second
)

type client struct {
	APIKey  string
	Client  *http.Client
	baseURL *url.URL
}

func NewClient(apiKey string) *client {
	parseUrl, _ := url.Parse(BaseURL)
	return &client{
		APIKey:  apiKey,
		Client:  httpClient(),
		baseURL: parseUrl,
	}
}

func httpClient() *http.Client {
	var transport http.RoundTripper = &http.Transport{
		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
		DisableKeepAlives:   false,
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 5 * time.Second,
		ForceAttemptHTTP2:   true,
	}

	client := &http.Client{
		Timeout:   DefaultHTTPTimeout,
		Transport: transport,
	}

	return client
}

func (c client) makeRequest(method string, path string, body any) ([]byte, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	parsedURL, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, parsedURL.String(), buf)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not execute request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body: %w", err)
	}

	return response, nil
}
