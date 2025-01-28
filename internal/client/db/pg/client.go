package pg

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/um3ra/auth-microservice/internal/client/db"
)

type pgClient struct {
	masterDbc db.DB
}

func NewClient(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return &pgClient{
		masterDbc: &pg{dbc: dbc},
	}, nil
}

func (c *pgClient) DB() db.DB {
	return c.masterDbc
}

func (c *pgClient) Close() error {
	if c.masterDbc != nil {
		c.masterDbc.Close()
	}
	return nil
}
