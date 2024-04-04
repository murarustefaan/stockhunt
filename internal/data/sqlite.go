package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDividendsStore struct {
	*sync.Mutex
	*sql.DB
	DividendsStore
}

func getDataPath() string {
	base := os.Getenv("SQLITE_VOLUME_PATH")
	if base == "" {
		return ":memory:"
	}

	return path.Join(base, "stockhunt.db")
}

func NewSqliteDividendsStore() (*SqliteDividendsStore, error) {
	datapath := getDataPath()
	db, err := sql.Open("sqlite3", datapath)
	if err != nil {
		return nil, err
	}

	store := SqliteDividendsStore{
		Mutex: &sync.Mutex{},
		DB:    db,
	}

	err = store.Migrate()
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *SqliteDividendsStore) Migrate() error {
	driver, err := sqlite3.WithInstance(s.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	s.Mutex.Lock()
	err = m.Up()
	s.Mutex.Unlock()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (s *SqliteDividendsStore) Add(data CompanyDividendInfo) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	_, err := s.Exec("INSERT INTO dividends (market, isin, name, value, yield, ex_date, pay_date, registration_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		data.Market, data.ISIN, data.Name, data.Value, data.Yield, data.ExDate, data.PaymentDate, data.RegistrationDate)
	return err
}

func (s *SqliteDividendsStore) List() ([]CompanyDividendInfo, error) {
	rows, err := s.Query("SELECT isin, name, value, yield, ex_date, pay_date, registration_date FROM dividends ORDER BY ex_date DESC, isin")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dividends []CompanyDividendInfo
	for rows.Next() {
		var dividend CompanyDividendInfo
		err = rows.Scan(&dividend.ISIN, &dividend.Name, &dividend.Value, &dividend.Yield, &dividend.ExDate, &dividend.PaymentDate, &dividend.RegistrationDate)
		if err != nil {
			return nil, err
		}

		dividends = append(dividends, dividend)
	}

	return dividends, nil
}
