package httpclient

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"
)

type (
	Client struct {
		baseUrl    string
		httpClient *http.Client

		timeout time.Duration
	}

	Response struct {
		res  *http.Response
		body []byte
	}
)

const DEFAULT_TIMEOUT = 10 * time.Second

// New func returns a Client struct
func New(baseUrl string, opts ...ClientOption) *Client {
	httpClient := &http.Client{Timeout: DEFAULT_TIMEOUT}
	client := &Client{httpClient: httpClient, baseUrl: baseUrl, timeout: DEFAULT_TIMEOUT}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

func (c Client) withBaseUrl(path string) string {
	return c.baseUrl + path
}

// TODO allow passing opts on the fly
func (c Client) Get(ctx context.Context, path string) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.withBaseUrl(path), nil)
	log.Printf("URL: %s", req.URL.String())
	if err != nil {
		return nil, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	res, err := c.httpClient.Do(req.WithContext(reqCtx))

	log.Printf("Context error: %s", reqCtx.Err())
	if err := reqCtx.Err(); err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	if err != nil {
		return nil, err
	}

	return &Response{res: res, body: body}, nil

}

// Stub RoundTripper in tests
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}
