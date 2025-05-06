package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// WithPool creates a new Queries instance with a connection pool
func WithPool(pool *pgxpool.Pool) *Queries {
	return &Queries{db: pool}
}
