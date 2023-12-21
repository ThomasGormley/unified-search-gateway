package httpclient

import (
	"net/http"
	"time"
)

type Client struct {
	baseUrl    string
	httpClient *http.Client

	headers map[string]Header
	query   map[string]string
	body    []byte
	timeout time.Duration
}

type Header struct {
	Value     string
	IsDefault bool
}

const DEFAULT_TIMEOUT = 10 * time.Second

// New func returns a Client struct
// func New(baseUrl string, opts ...ClientOption) *Client {
// 	httpClient := &http.Client{Timeout: DEFAULT_TIMEOUT}
// 	client := &Client{httpClient: httpClient, baseUrl: baseUrl, timeout: DEFAULT_TIMEOUT}

// 	for _, opt := range opts {
// 		opt(client)
// 	}

// 	return client
// }
