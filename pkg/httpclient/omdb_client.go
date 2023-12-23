package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/thomasgormley/unified-search-gateway/internal/configuration"
	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

type OmdbClient struct {
	HttpClient *http.Client
	BaseUrl    string
}

type Omdb struct {
	client *Client
}

func NewOmdb() *Omdb {

	httpClient := New("https://www.omdbapi.com")
	return &Omdb{client: httpClient}
}

type SearchResponse struct {
	Search []models.Omdb `json:"Search"`
}

func (oc *Omdb) Search(title string, contentType string, releaseYear string) ([]models.Omdb, error) {
	params := url.Values{
		"apikey": {configuration.Get().OmdbApiKey},
		"s":      {title},
		"type":   {contentType},
		"y":      {releaseYear},
	}
	reqUrl := fmt.Sprintf("?%s", params.Encode())
	ctx := context.Background()
	response, err := oc.client.Get(ctx, reqUrl)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}

	res := response.res

	switch res.StatusCode / 100 {
	case 4, 5:
		return nil, fmt.Errorf("HTTP error: %s", res.Status)
	}

	var omdbRespJson SearchResponse

	if err := json.Unmarshal(response.body, &omdbRespJson); err != nil {
		return nil, err
	}

	return omdbRespJson.Search, nil
}
