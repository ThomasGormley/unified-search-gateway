package search

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type FilterValidator interface {
	validate() error
}

type SearchFilters interface {
	FilterValidator
}

type SearchOptions struct {
	Query   string
	Page    int
	PerPage int
	Filters SearchFilters
}

// type Option = func(c *Customer)

// func NewSearchQueryOptions

func (so *SearchOptions) validate() error {
	// some pattern for collecting all the errors to send back to client as 400
	if err := so.Filters.validate(); err != nil {
		return err
	}

	if so.Page < 0 {
		return errors.New("page must be greater than or equal to 0")
	}

	if so.PerPage < 0 {
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

type SearchResult[T any] struct {
	Data T
	Err  error
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

func main(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	postFilters := &PostFilters{
		Author:      q.Get("author"),
		Topics:      strings.Split(q.Get("topic"), ","),
		PublishedAt: q.Get("publishedAt"),
		Label:       q.Get("label"),
		Has:         strings.Split(q.Get("has"), ","),
		Type:        q.Get("publishedAt"),
	}
	searchOptions := &SearchOptions{
		Query:   q.Get("q"),
		Page:    1,
		PerPage: 10,
		Filters: postFilters,
	}
	searchService := NewSearchPostService(searchOptions)

	// searchService := Search[[]models.Post]{
	// 	Options: searchOptions,
	// 	Queryer: &PostQueryer{},
	// }

	res, err := searchService.HandleSearch()

	if err != nil {
		panic("err")
	}

	fmt.Printf("Response data: %s", res)
}
