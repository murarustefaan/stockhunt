package markets

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	IndexISIN = iota
	IndexName
	IndexValue
	IndexYield
	IndexExDate
	IndexPaymentDate
	_
	IndexRegistrationDate
)

const (
	dateFormat = "2.01.2006"
)

func ParseDividendInfoPage(body io.Reader) ([]CompanyDividendInfo, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	tableRows := doc.Find("table.dataTable tbody tr")
	if tableRows == nil {
		return nil, fmt.Errorf("failed to find dividend info table")
	}

	results := make([]CompanyDividendInfo, 0)

	tableRows.Each(func(i int, row *goquery.Selection) {
		record := CompanyDividendInfo{}
		row.Find("td").Each(func(i int, cell *goquery.Selection) {
			switch i {
			case IndexISIN:
				record.ISIN = strings.TrimSpace(cell.Find("a strong").Text())
			case IndexName:
				record.Name = strings.TrimSpace(cell.Text())
			case IndexValue:
				raw := normalizeLocale(cell.Text())
				record.Value, _ = strconv.ParseFloat(raw, 64)
			case IndexYield:
				raw := normalizeLocale(cell.Text())
				record.Yield, _ = strconv.ParseFloat(raw, 64)
			case IndexExDate:
				record.ExDate, _ = time.Parse(dateFormat, strings.TrimSpace(cell.Text()))
			case IndexPaymentDate:
				record.PaymentDate, _ = time.Parse(dateFormat, strings.TrimSpace(cell.Text()))
			case IndexRegistrationDate:
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
