package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Response represents an HTTP response
type Response struct {
	Status     string
	StatusCode int
	Body       any
	Headers    http.Header
}

// Request makes an HTTP request and returns a Response
func Request(method, rawURL string, headers map[string]string, body io.Reader, skipSSL bool) (*Response, error) {
	var err error

	req, err := http.NewRequest(method, rawURL, body)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Create a custom http.Client with or without SSL verification
	client := &http.Client{}
	if skipSSL {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
	}, nil
}

func PostJson(rawURL string, headers map[string]string, body any, skipSSL bool) (*Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}
	return Request("POST", rawURL, headers, bytes.NewBuffer(reqBody), skipSSL)
}

// Get makes a GET request
func Get(rawURL string, headers map[string]string, skipSSL bool) (*Response, error) {
	return Request("GET", rawURL, headers, nil, skipSSL)
}

// Post makes a POST request
func Post(rawURL string, headers map[string]string, body io.Reader, skipSSL bool) (*Response, error) {
	return Request("POST", rawURL, headers, body, skipSSL)
}

func PostFile(rawURL string, headers map[string]string, filePath string, skipSSL bool) (*Response, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	headers["content-type"] = writer.FormDataContentType()
	return Request("POST", rawURL, headers, &body, skipSSL)
}

// AddQueryParams adds query parameters to a URL
func AddQueryParams(rawURL string, params map[string]string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	query := parsedURL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}
