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

func TestOmdbSearch(t *testing.T) {

	tests := []testCase{
		{
			name: "Errors on 4xx",
			mockResponse: &http.Response{
				StatusCode: 400,
				Body:       io.NopCloser(bytes.NewBufferString(`{"error": "not found"}`)),
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			mockClient := httpclient.OmdbClient{
				HttpClient: &http.Client{
					Transport: RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
						return tc.mockResponse, nil
					}),
				},
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
