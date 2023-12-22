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

type OmdbQueryer struct {
	SearchOptions[OmdbFilters]
}

func OmdbQuery(searchOptions SearchOptions[OmdbFilters]) (SearchItem, error) {
	omdbClient := httpclient.NewOmdb()
	resp, err := omdbClient.Search(searchOptions.Query, searchOptions.Filters.Type, searchOptions.Filters.Y)

	if err != nil {
		return ResultSet[models.Omdb]{}, err
	}

	return ResultSet[models.Omdb]{
		Data: resp,
		Type: "omdb",
	}, nil
}

func (o OmdbQueryer) Query() (SearchItem, error) {
	omdbClient := httpclient.NewOmdb()
	resp, err := omdbClient.Search(o.SearchOptions.Query, o.SearchOptions.Filters.Type, o.SearchOptions.Filters.Y)

	if err != nil {
		return ResultSet[models.Omdb]{}, err
	}

	return ResultSet[models.Omdb]{
		Data: resp,
		Type: "omdb",
	}, nil
}
