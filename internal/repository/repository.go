package repository

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Timeout time.Duration
}

type Repository struct {
	pool   *pgxpool.Pool
	config Config
}

func New(pool *pgxpool.Pool, config Config) *Repository {
	return &Repository{
		pool:   pool,
		config: config,
	}
}
