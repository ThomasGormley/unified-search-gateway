package search

import (
	"fmt"
	"log"

	"github.com/thomasgormley/unified-search-gateway/pkg/httpclient"
	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

// func NewOmdbFilters()

// S represents the movie title to search for.
// Type represents the content type to return.
// Y represents the year of release.
type OmdbFilters struct {
	Type string
	Y    string
}

func (f OmdbFilters) isValidContentType() error {
	log.Printf("filter.type: %s", f.Type)
	switch f.Type {
	case models.OmdbContentTypeMovie, models.OmdbContentTypeEpisode, models.OmdbContentTypeSeries, "":
		return nil
	default:
		return fmt.Errorf("filter type must be one of: %s, %s, %s", models.OmdbContentTypeMovie, models.OmdbContentTypeSeries, models.OmdbContentTypeEpisode)
	}
}

func (filters OmdbFilters) Validate() error {

	if err := filters.isValidContentType(); err != nil {
		return err
	}
	// validation logic for OmdbFilters
	return nil
}

func OmdbQuery(searchOptions SearchOptions[OmdbFilters]) SearchItem {
	omdbClient := httpclient.NewOmdb()
	resp, err := omdbClient.Search(searchOptions.Query, searchOptions.Filters.Type, searchOptions.Filters.Y)

	if err != nil {
		return SearchResult[models.Omdb]{
			Error: err.Error(),
			Type:  "omdb",
		}
	}

	return SearchResult[models.Omdb]{
		Data: resp,
		Type: "omdb",
	}
}
