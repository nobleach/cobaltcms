package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/nobleach/cobaltcms/internal/types"
)

type Storage interface {
	GetPublishedStatuses() ([]types.PublishedStatus, error)
}

type PostgresStore struct {
	pool    *pgxpool.Pool
	queries *Queries
}

func NewPostgresStore() (*PostgresStore, error) {
	// TODO get this into an env var or config
	connStr := "postgres://cobaltcms:cobaltcmspass@localhost:5432/cobaltcms?sslmode=disable"
	poolConfig, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to parse database URL")
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)

	queries := New(pool)

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to create connection pool.")
	}

	return &PostgresStore{
		pool:    pool,
		queries: queries,
	}, nil
}

var ErrInvalidInput = errors.New("Invalid Input")

func (s *PostgresStore) GetPublishedStatuses() ([]types.PublishedStatus, error) {
	ctx := context.Background()

	statusList, err := s.queries.ListAllPublishedStatuses(ctx)

	if err != nil {
		log.Fatal().Msgf("Could not fetch statuses list: %v", err)

		return nil, err
	}

	var statuses []types.PublishedStatus

	for _, element := range statusList {
		status := types.PublishedStatus{
			Id:     element.ID.String(),
			Status: element.Status,
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}
