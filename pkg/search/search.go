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
		Data []T    `json:"data"`
		Type string `json:"type"`
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

	// Queryer is an interface that represents a search query.
	// Implementations of this interface must provide a Query method that performs the search query and returns the result.
	Queryer[Return SearchItem] interface {
		Query() (Return, error)
	}
)

func (qr ResultSet[T]) GetType() string {
	return qr.Type
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

func HandleSearch[R SearchItem](qs ...Queryer[R]) ([]R, error) {
	resultsChan := make(chan *R)
	errChan := make(chan error)
	var wg sync.WaitGroup

	for _, queryer := range qs {
		wg.Add(1)
		go func(q Queryer[R]) {
			defer wg.Done()
			data, err := q.Query()
			if err != nil {
				errChan <- err
				resultsChan <- nil // Send nil for failed query
			} else {
				resultsChan <- &data
			}
		}(queryer)
	}

	// Close channels when all goroutines are done
	go func() {
		wg.Wait()
		close(resultsChan)
		close(errChan)
	}()

	var aggregatedData []R
	var errors []error
	for {
		select {
		case data, ok := <-resultsChan:
			if !ok {
				resultsChan = nil
			} else if data != nil {
				aggregatedData = append(aggregatedData, *data)
			}
		case err, ok := <-errChan:
			if !ok {
				errChan = nil
			} else {
				errors = append(errors, err)
			}
		}
		if resultsChan == nil && errChan == nil {
			break
		}
	}

	// Handle the case where all queries failed
	if len(aggregatedData) == 0 && len(errors) > 0 {
		return nil, errors[0] // or aggregate errors as needed
	}

	return aggregatedData, nil
}
