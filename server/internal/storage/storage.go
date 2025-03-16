package storage

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/nobleach/cobaltcms/internal/types"
)

type Storage interface {
	GetPublishedStatuses() ([]types.PublishedStatus, error)
	GetPublishedContentForDate(dateTime string) ([]ListPublishedContentForDateTimeRow, error)
	GetPublishedContentForId(id string, dateTime string) ([]GetPublishedContentByIdRow, error)
	SaveContent(newContent types.NewContent) (string, error)
}

type PostgresStore struct {
	pool    *pgxpool.Pool
	queries *Queries
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := os.Getenv("DATABASE_URL")
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

func (s *PostgresStore) GetPublishedContentForDate(dateTime string) ([]ListPublishedContentForDateTimeRow, error) {
	ctx := context.Background()

	publishedContentList, err := s.queries.ListPublishedContentForDateTime(ctx, dateTime)

	if err != nil {
		log.Fatal().Msgf("Could not fetch published content: %v", err)

		return nil, err
	}

	return publishedContentList, nil
}

func (s *PostgresStore) GetPublishedContentForId(id string, dateTime string) ([]GetPublishedContentByIdRow, error) {
	ctx := context.Background()

	pageIdUUID, err := uuid.Parse(id)
	if err != nil {
		log.Error().Msgf("Could not parse UUID from string")
		return nil, ErrInvalidInput
	}

	params := GetPublishedContentByIdParams{
		ToTimestamp:   dateTime,
		PageContentID: pageIdUUID,
	}

	publishedContentList, err := s.queries.GetPublishedContentById(ctx, params)

	if err != nil {
		log.Fatal().Msgf("Could not fetch published content: %v", err)

		return nil, err
	}

	return publishedContentList, nil
}

func (s *PostgresStore) SaveContent(newContent types.NewContent) (string, error) {
	// Validate input for SCHEDULED content
	if newContent.PublishedStatus == "SCHEDULED" {
		if newContent.PublishStartDateTime == "" || newContent.PublishEndDateTime == "" {
			return "", errors.New("SCHEDULED content must have both start and end date times")
		}
	}

	ctx := context.Background()
	dateFormat := "2021-11-22T12:34:56"
	startDate, err := time.Parse(dateFormat, newContent.PublishStartDateTime)
	endDate, err := time.Parse(dateFormat, newContent.PublishEndDateTime)

	// Prepare parameters for the query
	params := SaveNewContentParams{
		FragmentType:       newContent.FragmentType,
		Name:               newContent.Name,
		Body:               newContent.Body,
		ExtendedAttributes: newContent.ExtendedAttributes,
		PublishedStatus:    newContent.PublishedStatus,
		PublishStart:       pgtype.Timestamptz{Time: startDate, Valid: true},
		PublishEnd:         pgtype.Timestamptz{Time: endDate, Valid: true},
	}

	// Execute the query to save the content
	uuid, err := s.queries.SaveNewContent(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save content")
		return "", err
	}
	//
	// return contentID.String(), nil
	return uuid.String(), nil
}
