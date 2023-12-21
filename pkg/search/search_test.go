package search_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/thomasgormley/unified-search-gateway/pkg/search"
)

type TestFilters struct {
	shouldFail bool
}

var errTest = fmt.Errorf("fail in test")

func (tf TestFilters) Validate() error {

	if tf.shouldFail {
		return errTest
	}
	return nil
}

func TestNewSearchOptions(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		page    int
		perPage int
		filters TestFilters
		want    *search.SearchOptions[TestFilters]
		wantErr error
	}{
		{
			name:    "Default values",
			query:   "test",
			page:    -1,
			perPage: 0,
			filters: TestFilters{},
			want: &search.SearchOptions[TestFilters]{
				Query:   "test",
				Page:    0,
				PerPage: 10,
				Filters: TestFilters{},
			},
			wantErr: nil,
		},
		{
			name:    "Runs validation on filters",
			query:   "test",
			page:    0,
			perPage: 0,
			filters: TestFilters{
				shouldFail: true,
			},
			want:    nil,
			wantErr: errTest,
		},
		// Add more test cases here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := search.NewSearchOptions[TestFilters](tt.query, tt.page, tt.perPage, tt.filters)
			if err != tt.wantErr {
				t.Errorf("NewSearchOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSearchOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
