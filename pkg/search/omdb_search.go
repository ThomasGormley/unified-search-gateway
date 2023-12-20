package search

import (
	"fmt"

	"github.com/thomasgormley/unified-search-gateway/pkg/models"
)

// func NewOmdbFilters()

// S represents the movie title to search for.
// Type represents the content type to return.
// Y represents the year of release.
type OmdbFilters struct {
	S    string
	Type string
	Y    string
}

func (f OmdbFilters) isValidContentType() bool {
	switch f.Type {
	case models.OmdbContentTypeMovie, models.OmdbContentTypeEpisode, models.OmdbContentTypeSeries:
		return true
	default:
		return false
	}
}

func (filters OmdbFilters) validate() error {
	if !filters.isValidContentType() {
		return fmt.Errorf("filter type must be one of: %s, %s, %s", models.OmdbContentTypeMovie, models.OmdbContentTypeSeries, models.OmdbContentTypeEpisode)
	}

	// validation logic for OmdbFilters
	return nil
}

type OmdbQueryer struct{}

func (OmdbQueryer) Query(opts SearchOptions[OmdbFilters]) ([]models.Omdb, error) {
	// http request
	return make([]models.Omdb, 10), nil
}

func NewOmdbSearchService(opts SearchOptions[OmdbFilters]) *Search[[]models.Omdb, OmdbFilters] {
	return &Search[[]models.Omdb, OmdbFilters]{
		Options: opts,
		Queryer: OmdbQueryer{},
	}
}
