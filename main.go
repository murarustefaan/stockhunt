package main

import (
	"log"
	"os"

	"github.com/murarustefaan/stockhunt/cmd/scraper"
)

func main() {
	var worker scraper.Scraper
	var err error

	switch {
	case len(os.Args) > 1 && os.Args[1] == "bvb":
	default:
		worker, err = scraper.NewBvbScraper()
	}

	if err != nil {
		log.Fatalf("failed to create scraper: %v", err)
	}

	err = worker.Update()
	if err != nil {
		log.Fatalf("failed to update data: %v", err)
	}
}
