package requests

import (
	"bytes"
	"maps"
)

// Client giữ config chung
type Client struct {
	BaseURL string
	Headers map[string]string
	SkipSSL bool
}

// NewClient tạo Client mới
func NewClient(baseURL string, headers map[string]string, skipSSL bool) *Client {
	return &Client{
		BaseURL: baseURL,
		Headers: headers,
		SkipSSL: skipSSL,
	}
}

// Get gọi GET request
func (c *Client) Get(endpoint string, headers map[string]string) (*Response, error) {
	return Request("GET", c.BaseURL+endpoint, mergeHeaders(c.Headers, headers), nil, c.SkipSSL)
}

// Post gọi POST request
func (c *Client) Post(endpoint string, data []byte, headers map[string]string) (*Response, error) {
	return Request("POST", c.BaseURL+endpoint, mergeHeaders(c.Headers, headers), bytes.NewBuffer(data), c.SkipSSL)
}

// Put gọi PUT request
func (c *Client) Put(endpoint string, data []byte, headers map[string]string) (*Response, error) {
	return Request("PUT", c.BaseURL+endpoint, mergeHeaders(c.Headers, headers), bytes.NewBuffer(data), c.SkipSSL)
}

// Delete gọi DELETE request
func (c *Client) Delete(endpoint string, headers map[string]string) (*Response, error) {
	return Request("DELETE", c.BaseURL+endpoint, mergeHeaders(c.Headers, headers), nil, c.SkipSSL)
}

// helper gộp headers mặc định + headers truyền thêm
func mergeHeaders(defaults, overrides map[string]string) map[string]string {
	merged := map[string]string{}
	maps.Copy(merged, defaults)
	maps.Copy(merged, overrides)
	return merged
}
