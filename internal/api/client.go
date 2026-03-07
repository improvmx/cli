package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const BaseURL = "https://api.improvmx.com/v3"

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = strings.NewReader(string(data))
	}

	reqURL := BaseURL + path
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth("api", c.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "improvmx-cli/0.1.0")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr struct {
			Error   string `json:"error"`
			Errors  map[string][]string `json:"errors"`
			Success bool   `json:"success"`
		}
		if json.Unmarshal(respBody, &apiErr) == nil && apiErr.Error != "" {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, apiErr.Error)
		}
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) Get(path string) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, body)
}

func (c *Client) Put(path string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPut, path, body)
}

func (c *Client) Delete(path string) ([]byte, error) {
	return c.doRequest(http.MethodDelete, path, nil)
}

func QueryEncode(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	v := url.Values{}
	for key, val := range params {
		v.Set(key, val)
	}
	return "?" + v.Encode()
}
