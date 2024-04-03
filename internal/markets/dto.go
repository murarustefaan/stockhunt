package markets

import "time"

type CompanyDividendInfo struct {
	Name string
	ISIN string

	Yield float64
	Value float64

	ExDate           time.Time
	PaymentDate      time.Time
	RegistrationDate time.Time
}
