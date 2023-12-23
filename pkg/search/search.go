package search

import (
	"errors"
	"log"
	"sync"
)

type (
	// FilterCriteria is an interface that defines the behavior of filter criteria used in search operations.
	// Implementations of this interface must provide a Validate method that validates the filter criteria.
	FilterCriteria interface {
		Validate() error
	}

	// SearchResult is a generic struct that represents the result set of a search operation.
	// It contains a slice of data of type T and a string representing the type of the data.
	SearchResult[T any] struct {
		Data  []T    `json:"data"`
		Error string `json:"error"`
		Type  string `json:"type"`
	}

	// SearchItem is an interface that represents a resource that can be searched.
	// Implementations of this interface must provide a GetType method that returns the type of the resource.
	SearchItem interface {
		GetType() string
	}

	// SearchOptions is a struct that represents the options passed to a search function.
	// It contains the search query, page number, items per page, and filter criteria.
	SearchOptions[F FilterCriteria] struct {
		Query   string
		Page    int
		PerPage int
		Filters F
	}
)

func (qr SearchResult[T]) GetType() string {
	return qr.Type
}

func NewSearchOptions[F FilterCriteria](query string, page, perPage int, filters F) (*SearchOptions[F], error) {

	// Set default values for page and perPage if they are not within specific ranges
	if page < 0 {
		page = 0
	}
	if perPage <= 0 {
		perPage = 10
	}

	opts := &SearchOptions[F]{
		Query:   query,
		Page:    page,
		PerPage: perPage,
		Filters: filters,
	}

	err := opts.Validate()

	if err != nil {
		return nil, err
	}

	return opts, nil
}

type NoopFilter struct{}

func (noop NoopFilter) Validate() error {
	return nil
}

func (opts SearchOptions[F]) Validate() error {
	log.Printf("Validating options\n")
	// some pattern for collecting all the errors to send back to client as 400
	if err := opts.Filters.Validate(); err != nil {
		return err
	}
	if opts.Page < 0 {
		return errors.New("page must be greater than or equal to 0")
	}

	if opts.PerPage < 0 {
		return errors.New("page must be greater than 0")
	}

	return nil
}

type QueryFunc func() SearchItem

func HandleSearch(queryFns ...QueryFunc) []SearchItem {
	resultsChan := make(chan *SearchItem)
	var wg sync.WaitGroup

	for _, query := range queryFns {
		wg.Add(1)
		go func(q QueryFunc) {
			defer wg.Done()
			data := q()

			resultsChan <- &data
		}(query)
	}

	// Close channels when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	var aggregatedData []SearchItem
	for data := range resultsChan {
		if data != nil {
			aggregatedData = append(aggregatedData, *data)
		}
	}

	// Handle the case where all queries failed
	if len(aggregatedData) == 0 {
		return nil // or aggregate errors as needed
	}

	return aggregatedData
}

func QueryFnFrom[F FilterCriteria](fn func(opts SearchOptions[F]) SearchItem, opts SearchOptions[F]) func() SearchItem {
	return func() SearchItem {
		return fn(opts)
	}
}
