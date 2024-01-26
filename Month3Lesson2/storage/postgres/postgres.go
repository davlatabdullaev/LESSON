package postgres

import (
	"database/sql"
	"fmt"
	"mini_market/config"
	"mini_market/storage"
)

type Store struct {
	db *sql.DB
}

func New(cfg config.Config) (storage.IStorage, error) {
	url := fmt.Sprintf(`host = %s port = %s user = %s password = %s database = %s sslmode=disable`,
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return Store{}, err
	}

	return Store{
		db: db,
	}, nil
}

func (s Store) CloseDB() {
	s.db.Close()
}

func (s Store) Staff() storage.IStaffRepo {
	return NewStaffRepo(s.db)
}
