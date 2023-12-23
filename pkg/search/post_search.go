package search

import (
	"math/rand"
	"time"

	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

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
	// validation logic for PostFilters
	return nil
}

func PostQuery(searchOptions SearchOptions[PostFilters]) SearchItem {
	// http request
	d := time.Duration(rand.Intn(2))
	time.Sleep(d * time.Second)

	return SearchResult[models.Post]{
		Data: []models.Post{},
		Type: "post",
	}
}
