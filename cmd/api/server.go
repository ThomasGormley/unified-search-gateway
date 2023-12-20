package api

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/thomasgormley/unified-search-gateway/pkg/search"
)

func csv(s string) []string {
	return strings.Split(s, ",")
}

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/search", handleApi)

	addr := "localhost:8080" // e.g., "localhost:8080" for local development

	// Start the server
	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func handleApi(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling request", "URL", r.URL.Query())
	q := r.URL.Query()

	postFilters := &search.PostFilters{
		Author:      q.Get("author"),
		Topics:      csv(q.Get("topic")),
		PublishedAt: q.Get("publishedAt"),
		Label:       q.Get("label"),
		Has:         csv(q.Get("has")),
		Type:        q.Get("publishedAt"),
	}
	searchOptions := &search.SearchOptions{
		Query:   q.Get("q"),
		Page:    1,
		PerPage: 10,
		Filters: postFilters,
	}
	postSearch := search.NewPostSearchService(searchOptions)

	postSearchRes, err := postSearch.HandleSearch()

	if err != nil {
		log.Printf("Error: %+v", err)
	}

	postSearchJson, err := json.Marshal(postSearchRes)

	if err != nil {
		slog.Error("Error marshalling.", "err", err)
	}
	w.Write(postSearchJson)
}
