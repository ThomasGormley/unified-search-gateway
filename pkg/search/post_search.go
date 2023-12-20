package search

import (
	"fmt"

	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

// func NewPostFilters()

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

func (pq PostQueryer) Query(opts SearchOptions[PostFilters]) ([]models.Post, error) {
	// http request
	return make([]models.Post, 10), nil
}

func NewPostSearchService(opts SearchOptions[PostFilters]) *Search[[]models.Post, PostFilters] {
	return &Search[[]models.Post, PostFilters]{
		Options: opts,
		Queryer: PostQueryer{},
	}
}
