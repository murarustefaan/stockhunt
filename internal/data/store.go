package data

type Store interface {
}

type DividendsStore interface {
	Store
	Add(CompanyDividendInfo) error
	List() ([]CompanyDividendInfo, error)
}
