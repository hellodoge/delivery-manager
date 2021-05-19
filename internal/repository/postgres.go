package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	postgresQueriesFolder = "postgres"
)

func NewPostgresDB(cfg DatabaseConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

