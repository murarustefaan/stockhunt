package markets

import (
	"fmt"
	"github.com/murarustefaan/stockhunt/internal/data"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	indexISIN = iota
	indexName
	indexValue
	indexYield
	indexExDate
	indexPaymentDate
	_
	indexRegistrationDate
)

const (
	dateFormat = "2.01.2006"
)

func ParseDividendInfoPage(body io.Reader) ([]data.CompanyDividendInfo, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	tableRows := doc.Find("table.dataTable tbody tr")
	if tableRows == nil {
		return nil, fmt.Errorf("failed to find dividend info table")
	}

	results := make([]data.CompanyDividendInfo, 0)

	tableRows.Each(func(i int, row *goquery.Selection) {
		record := data.CompanyDividendInfo{
			Market: "BVB",
		}
		row.Find("td").Each(func(i int, cell *goquery.Selection) {
			switch i {
			case indexISIN:
				record.ISIN = strings.TrimSpace(cell.Find("a strong").Text())
			case indexName:
				record.Name = strings.TrimSpace(cell.Text())
			case indexValue:
				raw := normalizeLocale(cell.Text())
				record.Value, _ = strconv.ParseFloat(raw, 64)
			case indexYield:
				raw := normalizeLocale(cell.Text())
				record.Yield, _ = strconv.ParseFloat(raw, 64)
			case indexExDate:
				record.ExDate, _ = time.Parse(dateFormat, strings.TrimSpace(cell.Text()))
			case indexPaymentDate:
				record.PaymentDate, _ = time.Parse(dateFormat, strings.TrimSpace(cell.Text()))
			case indexRegistrationDate:
				record.RegistrationDate, _ = time.Parse(dateFormat, strings.TrimSpace(cell.Text()))
			}
		})

		results = append(results, record)
	})

	return results, nil
}

func normalizeLocale(s string) string {
	return strings.Replace(
		strings.Replace(s, ".", "", -1),
		",", ".", -1,
	)
}
