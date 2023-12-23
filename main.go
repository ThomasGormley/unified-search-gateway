package main

import (
	"log"

	"github.com/thomasgormley/unified-search-gateway/cmd/api"
)

func main() {
	if e := api.Start(); e != nil {
		log.Fatalf("failed to start server: %v", e)
	}
}
