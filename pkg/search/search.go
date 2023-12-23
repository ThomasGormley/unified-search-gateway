package search

import (
	"errors"
	"log"
	"strconv"
	"sync"
)

type (
	// FilterCriteria is an interface that defines the behavior of filter criteria used in search operations.
	// Implementations of this interface must provide a Validate method that validates the filter criteria.
	FilterCriteria interface {
		Validate() error
	}

	// ResultSet is a generic struct that represents the result set of a search operation.
	// It contains a slice of data of type T and a string representing the type of the data.
	ResultSet[T any] struct {
		Data  []T     `json:"data"`
		Error *string `json:"error"`
		Type  string  `json:"type"`
	}

	// SearchItem is an interface that represents a resource that can be searched.
	// Implementations of this interface must provide a GetType method that returns the type of the resource.
	SearchItem interface {
		GetType() string
		GetError() *string
	}

	// SearchOptions is a struct that represents the options passed to a search function.
	// It contains the search query, page number, items per page, and filter criteria.
	SearchOptions[F FilterCriteria] struct {
		Query   string
		Page    int
		PerPage int
		Filters F
	}

	// Queryer is an interface that represents a search query.
	// Implementations of this interface must provide a Query method that performs the search query and returns the result.
	Queryer interface {
		Query() SearchItem
	}
)

func (qr ResultSet[T]) GetType() string {
	return qr.Type
}

func (qr ResultSet[T]) GetError() *string {
	log.Printf("GetError: %+v", qr.Error)
	if qr.Error == nil {
		log.Printf("Returning nil")
		return nil
	}
	return qr.Error
}

func NewSearchOptions[F FilterCriteria](query string, p, pPage string, filters F) (*SearchOptions[F], error) {
	// convert page and perPage to int
	page := 0
	if p != "" {
		var err error
		page, err = strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
	}

	perPage := 10
	if pPage != "" {
		var err error
		perPage, err = strconv.Atoi(pPage)
		if err != nil {
			return nil, err
		}
	}

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

func HandleSearch(qs ...QueryFunc) []SearchItem {
	resultsChan := make(chan *SearchItem)
	var wg sync.WaitGroup

	for _, queryer := range qs {
		wg.Add(1)
		go func(q QueryFunc) {
			defer wg.Done()
			data := q()

			resultsChan <- &data
		}(queryer)
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
