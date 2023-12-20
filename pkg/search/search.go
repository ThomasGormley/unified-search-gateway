package search

import (
	"errors"
	"log/slog"
)

// SearchFilters is an interface that represents the filters used for searching.
type SearchFilters interface {
	validate() error
}

type SearchOptions struct {
	Query   string
	Page    int
	PerPage int
	Filters SearchFilters
}

func (opts *SearchOptions) validate() error {
	// some pattern for collecting all the errors to send back to client as 400
	if err := opts.Filters.validate(); err != nil {
		slog.Info("error validating")
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

type Queryer[T any] interface {
	Query(*SearchOptions) (T, error)
}
type Search[T any] struct {
	Queryer[T]
	Options *SearchOptions
}

func (s *Search[T]) HandleSearch() (*T, error) {

	if err := s.Options.validate(); err != nil {
		return nil, err
	}

	resultData, err := s.Queryer.Query(s.Options)
	if err != nil {
		return nil, err
	}

	return &resultData, nil
}
