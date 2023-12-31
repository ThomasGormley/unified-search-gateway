package api

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strings"

	"github.com/thomasgormley/unified-search-gateway/internal/configuration"
	"github.com/thomasgormley/unified-search-gateway/pkg/search"
)

func csv(s string) []string {
	return strings.Split(s, ",")
}

func Start() {
	configuration.Load()
	mux := http.NewServeMux()

	// mux.HandleFunc("/api/search", handleApi)
	mux.HandleFunc("/api/search/omdb", handleOmdbSearch)

	addr := "localhost:8080" // e.g., "localhost:8080" for local development

	// Start the server
	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendError(w http.ResponseWriter, errorCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	resp := ErrorResponse{Error: message}
	json.NewEncoder(w).Encode(resp)
}

func sendInternalServerError(w http.ResponseWriter) {
	sendError(w, http.StatusInternalServerError, "Something went wrong.")
}

func handleOmdbSearch(w http.ResponseWriter, r *http.Request) {
	slog.Info("Handling OMDB Search", "Query", r.URL.Query())
	q := r.URL.Query()
	omdbFilters := search.OmdbFilters{
		Type: q.Get("type"),
		Y:    q.Get("year"),
	}

	omdbSearchOpts, err := search.NewSearchOptions[search.OmdbFilters](q.Get("q"), 1, 10, omdbFilters)

	if err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	omdbSearch := search.NewOmdbSearchService(*omdbSearchOpts)

	omdbSearchRes, err := omdbSearch.HandleSearch()

	if err != nil {
		log.Printf("Error: %+v", err)
		sendInternalServerError(w)
		return
	}

	postSearchJson, err := json.Marshal(omdbSearchRes)

	if err != nil {
		slog.Error("Error marshalling.", "err", err)
		sendInternalServerError(w)
		return
	}

	w.Write(postSearchJson)
}

// func handleApi(w http.ResponseWriter, r *http.Request) {
// 	slog.Info("Handling request", "URL", r.URL.Query())
// 	q := r.URL.Query()

// 	postFilters := search.PostFilters{
// 		Author:      q.Get("author"),
// 		Topics:      csv(q.Get("topic")),
// 		PublishedAt: q.Get("publishedAt"),
// 		Label:       q.Get("label"),
// 		Has:         csv(q.Get("has")),
// 		Type:        q.Get("publishedAt"),
// 	}
// 	searchOptions := search.SearchOptions[search.PostFilters]{
// 		Query:   q.Get("q"),
// 		Page:    1,
// 		PerPage: 10,
// 		Filters: postFilters,
// 	}
// 	postSearch := search.NewPostSearchService(searchOptions)

// 	postSearchRes, err := postSearch.HandleSearch()

// 	if err != nil {
// 		log.Printf("Error: %+v", err)
// 	}

// 	postSearchJson, err := json.Marshal(postSearchRes)

// 	if err != nil {
// 		slog.Error("Error marshalling.", "err", err)
// 	}
// 	w.Write(postSearchJson)
// }
