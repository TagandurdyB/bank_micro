package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (c *PgxClient) Exec(sql string, args ...any) (pgconn.CommandTag, error) {
	return c.Conn.Exec(c.Ctx, sql, args...)
}
