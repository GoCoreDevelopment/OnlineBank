package db

import (
	"api/internal/config"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

func InitDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DataBaseURL)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(1 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(5)

	err = db.Ping()
	if err != nil {
		return nil, errors.New("DB connection fail")
	}

	return db, nil
}

func CloseDB(db *sql.DB) error{
	err := db.Close()
	if err != nil {
		return err
	}

	return nil
}
