package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func LiveDbConn() (*pgx.Conn, error) {
	live_db_dsn := os.Getenv("LIVE_DB_DSN")
	conn, err := pgx.Connect(context.Background(), live_db_dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func ViewDbConn() (*pgx.Conn, error) {
	view_db_dsn := os.Getenv("VIEW_DB_DSN")
	conn, err := pgx.Connect(context.Background(), view_db_dsn)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
