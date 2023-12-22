package httpclient_test

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/thomasgormley/unified-search-gateway/pkg/httpclient"
	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

// For stubbing
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

var mockOmdbClient = httpclient.OmdbClient{
	HttpClient: &http.Client{
		Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
			// Assert on request attributes
			// Return a response or error you want
			return &http.Response{}, nil
		}),
	},
}

type testCase struct {
	name         string
	title        string
	contentType  string
	releaseYear  string
	mockResponse *http.Response // Mock response to return from the RoundTripperFunc
	expected     []models.Omdb  // Expected result from the Search method
	expectError  bool           // Whether an error is expected
}

var (
	errorResponse       = `{"error": "not found"}`
	successResponseJson = `{
						"Search": [
						{
							"Title": "The Good, the Bad and the Ugly",
							"Year": "1966",
							"Rated": "",
							"Released": "",
							"Runtime": "",
							"Genre": "",
							"Director": "",
							"Writer": "",
							"Actors": "",
							"Plot": "",
							"Language": "",
							"Country": "",
							"Awards": "",
							"Poster": "https://m.media-amazon.com/images/M/MV5BNjJlYmNkZGItM2NhYy00MjlmLTk5NmQtNjg1NmM2ODU4OTMwXkEyXkFqcGdeQXVyMjUzOTY1NTc@._V1_SX300.jpg",
							"Ratings": null,
							"Metascore": "",
							"imdbRating": "",
							"imdbVotes": "",
							"imdbID": "tt0060196",
							"Type": "movie",
							"DVD": "",
							"BoxOffice": "",
							"Production": "",
							"Website": "",
							"Response": ""
						}]
						}`
)

func TestOmdbSearch(t *testing.T) {

	tests := []testCase{
		{
			name: "returns search result",
			mockResponse: &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(successResponseJson)),
			},
			expected: []models.Omdb{{
				Title:  "The Good, the Bad and the Ugly",
				Year:   "1966",
				ImdbID: "tt0060196",
				Type:   "movie",
				Poster: "https://m.media-amazon.com/images/M/MV5BNjJlYmNkZGItM2NhYy00MjlmLTk5NmQtNjg1NmM2ODU4OTMwXkEyXkFqcGdeQXVyMjUzOTY1NTc@._V1_SX300.jpg",
			}},
			expectError: false,
		},
		{
			name: "Errors on 4xx",
			mockResponse: &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewBufferString(errorResponse)),
			},
			expected:    nil,
			expectError: true,
		},
		{
			name: "Errors on 5xx",
			mockResponse: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(bytes.NewBufferString(errorResponse)),
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mockClient := httpclient.Omdb{
				Client: httpclient.New("https://www.omdbapi.com", httpclient.WithCustomHttpClient(&http.Client{
					Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
						return tc.mockResponse, nil
					}),
				})),
			}

			result, err := mockClient.Search(tc.title, tc.contentType, tc.releaseYear)

			if (err != nil) != tc.expectError {
				t.Errorf("Expected error: %v, got: %v", tc.expectError, err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected result: %+v, got %+v", tc.expected, result)
			}

		})
	}
}
