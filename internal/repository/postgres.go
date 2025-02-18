package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(databaseURL string) (*PostgresRepository, error) {
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	// Testa a conexão
	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r *PostgresRepository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}
