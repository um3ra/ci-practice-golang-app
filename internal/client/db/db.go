package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Handler is a function that is executed within a transaction
type Handler func(ctx context.Context) error

// Client represents a database client
type Client interface {
	DB() DB
	Close() error
}

// TxManager is responsible for handling transactions
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query is a wrapper for a query containing the query name and the raw query
// The query name is used for logging and potentially could be used elsewhere, e.g., for tracing
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor interface for handling database transactions
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

// SQLExecer combines NamedExecer and QueryExecer interfaces
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer interface for working with named queries using tags in structures
type NamedExecer interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecer interface for working with regular queries
type QueryExecer interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger interface for checking database connection health
type Pinger interface {
	Ping(ctx context.Context) error
}

// DB interface represents the overall database client, combining SQLExecer, Transactor, and Pinger interfaces
type DB interface {
	SQLExecer
	Transactor
	Pinger
	Close()
}
