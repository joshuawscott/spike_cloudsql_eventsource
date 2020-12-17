package database

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v4"
)

// CreateTable creates a table from a CREATE TABLE statement query
func CreateTable(conn *pgx.Conn, query string) error {
	res, err := conn.Exec(context.Background(), query)
	if err != nil {
		return err
	}
	num := res.RowsAffected()
	if num != 0 {
		return errors.New("Got non-zero rows altered in response")
	}
	return nil
}

// CreateConn reads the PGURL environment variable and connects to that database.
func CreateConn() (*pgx.Conn, error) {
	postgresURL := os.Getenv("PGURL")

	return pgx.Connect(context.Background(), postgresURL)
}
