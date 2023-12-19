package search

import (
	"fmt"

	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

func NewPostFilters()

type PostFilters struct {
	Author      string
	Topics      []string
	PublishedAt string
	Label       string
	// has video, image, link, note within body
	Has []string
	// content type as published
	Type string
}

func (pf PostFilters) validate() error {
	if pf.Author == "" {
		return fmt.Errorf("author cannot be empty")
	}
	// validation logic for PostFilters
	return nil
}

type PostQueryer struct{}

func (pq *PostQueryer) Query(*SearchOptions) ([]models.Post, error) {
	// http request
	return make([]models.Post, 10), nil
}

func NewSearchPostService(opts *SearchOptions) *Search[[]models.Post] {
	return &Search[[]models.Post]{
		Options: opts,
		Queryer: &PostQueryer{},
	}
}
