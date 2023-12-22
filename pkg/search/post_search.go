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

func (pf PostFilters) Validate() error {
	if pf.Author == "" {
		return fmt.Errorf("author cannot be empty")
	}
	// validation logic for PostFilters
	return nil
}

type PostQueryer struct {
	SearchOptions[PostFilters]
}

func (pq PostQueryer) Query() (SearchableResource, error) {
	// http request
	return QueryResult[models.Post]{
		Data: []models.Post{},
		Type: "post",
	}, nil

}
