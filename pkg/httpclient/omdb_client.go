package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/thomasgormley/unified-search-gateway/internal/configuration"
	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

type OmdbClient struct {
	HttpClient *http.Client
	BaseUrl    string
}

// For stubbing
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

var mockOmdbClient = OmdbClient{
	HttpClient: &http.Client{
		Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			// Assert on request attributes
			// Return a response or error you want
			return &http.Response{}, nil
		}),
	},
}

func NewOmdb() OmdbClient {

	httpClient := &http.Client{Timeout: DEFAULT_TIMEOUT}
	return OmdbClient{
		HttpClient: httpClient,
		BaseUrl:    "https://www.omdbapi.com",
	}
}

type searchResponse struct {
	Search []models.Omdb `json:"Search"`
}

func (client *OmdbClient) Search(title string, contentType string, releaseYear string) ([]models.Omdb, error) {
	params := url.Values{
		"apikey": {configuration.Get().OmdbApiKey},
		"s":      {title},
		"type":   {contentType},
		"y":      {releaseYear},
	}
	reqUrl := fmt.Sprintf("%s?%s", client.BaseUrl, params.Encode())
	slog.Info("OMDB API", "URL", reqUrl)
	req, err := http.NewRequest(http.MethodGet, reqUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.HttpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	fmt.Printf("status: %d, mod: %d", resp.StatusCode, resp.StatusCode/100)
	switch resp.StatusCode / 100 {
	case 4, 5:
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	var omdbRespJson searchResponse

	if err := json.Unmarshal(body, &omdbRespJson); err != nil {
		return nil, err
	}

	return omdbRespJson.Search, nil
}
