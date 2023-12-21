package search

import (
	"errors"
	"fmt"
	"log"
)

// SearchFilters is an interface that represents the filters used for searching.
type SearchFilters interface {
	validate() error
}

type SearchOptions[F SearchFilters] struct {
	Query   string
	Page    int
	PerPage int
	Filters F
}

func NewSearchOptions[F SearchFilters](query string, page, perPage int, filters SearchFilters) (*SearchOptions[F], error) {
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
		Filters: filters.(F),
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
	if err := opts.Filters.validate(); err != nil {
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

type Queryer[R any, F SearchFilters] interface {
	Query(SearchOptions[F]) (R, error)
}
type Search[R any, F SearchFilters] struct {
	Queryer[R, F]
	Options SearchOptions[F]
}

func (s Search[R, F]) HandleSearch() (*R, error) {

	resultData, err := s.Queryer.Query(s.Options)

	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("err")

	return &resultData, nil
}
