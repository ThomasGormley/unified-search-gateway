package search

import (
	"fmt"

	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

func NewPeopleFilters()

type PeopleFilters struct {
	Author      string
	Topics      []string
	PublishedAt string
	Label       string
	// has video, image, link, note within body
	Has []string
	// content type as published
	Type string
}

func (pf PeopleFilters) validate() error {
	if pf.Author == "" {
		return fmt.Errorf("author cannot be empty")
	}
	// validation logic for PeopleFilters
	return nil
}

type PeopleQueryer struct{}

func (pq *PeopleQueryer) Query(*SearchOptions) ([]models.People, error) {
	// http request
	return make([]models.People, 10), nil
}

func NewSearchPeopleService(opts *SearchOptions) *Search[[]models.People] {
	return &Search[[]models.People]{
		Options: opts,
		Queryer: &PeopleQueryer{},
	}
}
