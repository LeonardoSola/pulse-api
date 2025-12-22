package postgres

import (
	"database/sql"
	"fmt"
)

func NewPostgres() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
}
