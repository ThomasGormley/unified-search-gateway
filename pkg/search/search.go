package search

import (
	"errors"
	"log"
)

type (
	// SearchFilters is an interface that represents the filters used for searching.
	SearchFilters interface {
		Validate() error
	}

	QueryResult[T any] struct {
		Data []T    `json:"data"`
		Type string `json:"type"`
	}

	SearchableResource interface {
		GetType() string
	}

	SearchOptions[F SearchFilters] struct {
		Query   string
		Page    int
		PerPage int
		Filters F
	}

	Queryer[QueryR SearchableResource] interface {
		Query() (QueryR, error)
	}
	Search[R SearchableResource] struct {
		Queriers []Queryer[R]
	}
)

func (qr QueryResult[T]) GetType() string {
	return qr.Type
}

func NewSearchOptions[F SearchFilters](query string, page, perPage int, filters F) (*SearchOptions[F], error) {
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

func NewSearch[R SearchableResource](queriers ...Queryer[R]) *Search[R] {
	return &Search[R]{
		Queriers: queriers,
	}
}

func (s Search[R]) HandleSearch() ([]R, error) {
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
