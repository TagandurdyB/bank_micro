package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PgxClient struct {
	Conn *pgx.Conn
	Ctx  context.Context
}

func ConnectPgx(dbURL string) (*PgxClient, error) {
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return &PgxClient{Conn: conn, Ctx: context.Background()}, nil
}

func (c *PgxClient) QueryRow(sql string, args ...any) pgx.Row {
	return c.Conn.QueryRow(c.Ctx, sql, args...)
}

func (c *PgxClient) Query(sql string, args ...any) (pgx.Rows, error) {
	return c.Conn.Query(c.Ctx, sql, args...)
}
