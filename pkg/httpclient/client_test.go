package httpclient

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func Test_New_ReturnsNewClient(t *testing.T) {

	tests := []struct {
		name string
		opts []ClientOption
		want *Client
	}{
		{
			name: "with options",
			opts: []ClientOption{
				WithTimeout(10 * time.Second),
				WithCustomHttpClient(&http.Client{}),
			},
			want: &Client{
				httpClient: &http.Client{},
				timeout:    10 * time.Second,
			},
		},
		{
			name: "with default timeout",
			opts: []ClientOption{},
			want: &Client{
				httpClient: &http.Client{Timeout: DEFAULT_TIMEOUT},
				timeout:    DEFAULT_TIMEOUT,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := New("", tc.opts...)

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Expected result: %+v, got %+v", tc.want, got)
			}
		})
	}

}
