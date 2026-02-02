package repository

import "database/sql"

type sqliteRepo struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *sqliteRepo {
	return &sqliteRepo{db: db}
}
