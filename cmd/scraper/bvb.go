package scraper

import (
	"context"
	"fmt"
	"github.com/murarustefaan/stockhunt/internal/data"
	"net/http"
	"slices"
	"time"

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

func (s *BvbScraper) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

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

	results, err := markets.ParseDividendInfoPage(res.Body)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		return fmt.Errorf("no results found")
	}

	existing, err := s.store.List()
	if err != nil {
		return fmt.Errorf("failed to list existing dividend info: %w", err)
	}

	for _, result := range results {
		if slices.ContainsFunc(existing, func(i data.CompanyDividendInfo) bool {
			return i.ISIN == result.ISIN && i.ExDate == result.ExDate
		}) {
			continue
		}

		err = s.store.Add(result)
		if err != nil {
			return fmt.Errorf("failed to save dividend info: %w", err)
		}
	}

	return nil
}
