package markets

import (
	"os"
	"path"
	"testing"
)

func TestParseDividendInfoPage(t *testing.T) {
	wd, _ := os.Getwd()

	data, err := os.Open(path.Join(wd, "../../testdata/bvb_dividend_info.html"))
	if err != nil {
		t.Fatalf("failed to open test data: %v", err)
	}

	results, err := ParseDividendInfoPage(data)
	if err != nil {
		t.Fatalf("failed to parse dividend info page: %v", err)
	}

	if len(results) == 0 {
		t.Fatalf("expected non-empty results")
	}
}
