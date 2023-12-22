package search

import (
	"errors"
	"log"
)

type ( // FilterCriteria is an interface that defines the behavior of filter criteria used in search operations.
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

	// Identifiable is an interface that represents a resource that can be searched.
	// Implementations of this interface must provide a GetType method that returns the type of the resource.
	Identifiable interface {
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
	Queryer[QueryR Identifiable] interface {
		Query() (QueryR, error)
	}

	// SearchService is a struct that represents a search service.
	// It contains a list of queryers that can be used to perform search queries.
	SearchService[R Identifiable] struct {
		Queriers []Queryer[R]
	}
)

func (qr ResultSet[T]) GetType() string {
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

func NewSearchService[R Identifiable](queriers ...Queryer[R]) *SearchService[R] {
	return &SearchService[R]{
		Queriers: queriers,
	}
}

func (s SearchService[R]) HandleSearch() ([]R, error) {
	var aggregatedData []R
	for _, queryer := range s.Queriers {
		data, err := queryer.Query()
		if err != nil {
			return nil, err
		}
		aggregatedData = append(aggregatedData, data)
	}

	return aggregatedData, nil
}
