package service

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type FSClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewFSClient(baseURL, apiKey string) *FSClient {
	return &FSClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *FSClient) Do(method, path string, uid, gid int, groups string, body io.Reader, contentType string, extraHeaders map[string]string) (*http.Response, error) {
	url := c.BaseURL + path

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-FS-UID", fmt.Sprintf("%d", uid))
	req.Header.Set("X-FS-GID", fmt.Sprintf("%d", gid))
	req.Header.Set("X-FS-Groups", groups)

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}

	return c.HTTPClient.Do(req)
}

// DoStream performs a request using an HTTP client without timeout, suitable for SSE streams.
func (c *FSClient) DoStream(method, path string, uid, gid int, groups string) (*http.Response, error) {
	url := c.BaseURL + path

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("X-FS-UID", fmt.Sprintf("%d", uid))
	req.Header.Set("X-FS-GID", fmt.Sprintf("%d", gid))
	req.Header.Set("X-FS-Groups", groups)

	// Use a client without timeout for long-lived SSE connections
	streamClient := &http.Client{}
	return streamClient.Do(req)
}
