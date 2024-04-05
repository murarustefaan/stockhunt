package scraper

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/murarustefaan/stockhunt/internal/data"
	"github.com/murarustefaan/stockhunt/internal/markets"
)

const (
	endpoint        = "https://bvb.ro/FinancialInstruments/CorporateActions/InfoDividend"
	userAgentHeader = "Stock Notifications Project (mailto: muraru.stefaan@gmail.com)"
)

type BvbScraper struct {
	endpoint string
	store    data.DividendsStore
}

func NewBvbScraper() (*BvbScraper, error) {
	store, err := data.NewSqliteDividendsStore()
	if err != nil {
		return nil, err
	}

	return &BvbScraper{
		endpoint: endpoint,
		store:    store,
	}, nil
}

func (s *BvbScraper) Update() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Print("Updating BVB dividend info data...")

	start := time.Now()
	ticker := time.Now()
	defer func() {
		if err != nil {
			log.Printf("Failed to update BVB dividend info data: %v", err)
		} else {
			log.Printf("Update BVB dividend info data done, took %s", time.Since(start))
		}
	}()

	log.Printf("Fetching BVB dividend info from %s", s.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgentHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	log.Printf("Fetched BVB dividend info data, took %s", time.Since(ticker))
	ticker = time.Now()

	log.Printf("Parsing BVB dividend info data...")
	results, err := markets.ParseDividendInfoPage(res.Body)
	if err != nil {
		return err
	}

	log.Printf("Parsed BVB dividend info data, took %s", time.Since(ticker))

	if len(results) == 0 {
		return fmt.Errorf("no results found")
	}

	existing, err := s.store.List()
	if err != nil {
		return fmt.Errorf("failed to list existing dividend info: %w", err)
	}

	ticker = time.Now()
	log.Printf("Saving BVB dividend info data...")

	for _, result := range results {
		if slices.ContainsFunc(existing, func(i data.CompanyDividendInfo) bool {
			return i.ISIN == result.ISIN && i.ExDate == result.ExDate
		}) {
			log.Printf("Skipping duplicate dividend info for %s", result.ISIN)
			continue
		}

		log.Printf("Saving dividend info for %s", result.ISIN)
		err = s.store.Add(result)
		if err != nil {
			return fmt.Errorf("failed to save dividend info: %w", err)
		}
	}

	log.Printf("Saved BVB dividend info data, took %s", time.Since(ticker))
	return nil
}
