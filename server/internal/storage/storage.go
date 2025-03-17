package storage

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/nobleach/cobaltcms/internal/types"
)

type Storage interface {
	GetPublishedStatuses() ([]types.PublishedStatus, error)
	GetPublishedContentForDate(dateTime string) ([]ListPublishedContentForDateTimeRow, error)
	GetPublishedContentForId(id string, dateTime string) ([]GetPublishedContentByIdRow, error)
	SaveContent(newContent types.NewContent) (string, error)
	UpdateContent(updateContent types.UpdateContent) (types.UpdateContent, error)
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
		log.Info().Msg("Fragment is SCHEDULED")
		if newContent.PublishStartDateTime == "" || newContent.PublishEndDateTime == "" {
			return "", errors.New("SCHEDULED content must have both start and end date times")
		}
	}

	ctx := context.Background()

	// Prepare parameters for the query
	params := SaveNewContentParams{
		FragmentType:       newContent.FragmentType,
		Name:               newContent.Name,
		Body:               newContent.Body,
		ExtendedAttributes: newContent.ExtendedAttributes,
		PublishedStatus:    newContent.PublishedStatus,
		PublishStart:       nil,
		PublishEnd:         nil,
	}

	if newContent.PublishedStatus == "SCHEDULED" {
		dateFormat := "2006-01-02 15:04:05"
		startDate, err := time.Parse(dateFormat, newContent.PublishStartDateTime)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse start date time")
			return "", errors.New("invalid start date time format")
		}

		endDate, err := time.Parse(dateFormat, newContent.PublishEndDateTime)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse end date time")
			return "", errors.New("invalid end date time format")
		}

		log.Debug().Msgf("Parsing start time %v as a timestamp", newContent.PublishStartDateTime)
		log.Debug().Msgf("Parsing end time %v as a timestamp", newContent.PublishEndDateTime)

		params.PublishStart = &startDate
		params.PublishEnd = &endDate
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

func (s *PostgresStore) UpdateContent(updateContent types.UpdateContent) (types.UpdateContent, error) {
	// Validate input for SCHEDULED content
	if updateContent.PublishedStatus == "SCHEDULED" {
		log.Info().Msg("Fragment is SCHEDULED")
		if updateContent.PublishStartDateTime == "" || updateContent.PublishEndDateTime == "" {
			return types.UpdateContent{}, errors.New("SCHEDULED content must have both start and end date times")
		}
	}

	ctx := context.Background()

	id, err := uuid.Parse(updateContent.Id)
	if err != nil {
		return types.UpdateContent{}, errors.New("Could not parse id into a UUID")
	}

	// Prepare parameters for the query
	params := UpdateContentParams{
		ID:                 id,
		FragmentType:       updateContent.FragmentType,
		Name:               updateContent.Name,
		Body:               updateContent.Body,
		ExtendedAttributes: updateContent.ExtendedAttributes,
		PublishedStatus:    updateContent.PublishedStatus,
		PublishStart:       nil,
		PublishEnd:         nil,
	}

	if updateContent.PublishedStatus == "SCHEDULED" {
		dateFormat := "2006-01-02 15:04:05"
		startDate, err := time.Parse(dateFormat, updateContent.PublishStartDateTime)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse start date time")
			return types.UpdateContent{}, errors.New("Failed to parse time into valid timestamp")
		}

		endDate, err := time.Parse(dateFormat, updateContent.PublishEndDateTime)
		if err != nil {
			log.Error().Err(err).Msg("Failed to parse end date time")
			return types.UpdateContent{}, errors.New("Failed to parse time into valid timestamp")
		}

		log.Debug().Msgf("Parsing start time %v as a timestamp", updateContent.PublishStartDateTime)
		log.Debug().Msgf("Parsing end time %v as a timestamp", updateContent.PublishEndDateTime)

		params.PublishStart = &startDate
		params.PublishEnd = &endDate
	}

	// Execute the query to save the content
	record, err := s.queries.UpdateContent(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save content")
		return types.UpdateContent{}, err
	}

	result := types.UpdateContent{
		Id:                   record.ID.String(),
		FragmentType:         record.FragmentType,
		Name:                 record.Name,
		Body:                 record.Body,
		ExtendedAttributes:   record.ExtendedAttributes,
		PublishedStatus:      record.PublishedStatus,
		PublishStartDateTime: "",
		PublishEndDateTime:   "",
	}

	if record.PublishStart != nil {
		result.PublishStartDateTime = record.PublishStart.String()
	}

	if record.PublishEnd != nil {
		result.PublishEndDateTime = record.PublishEnd.String()
	}

	return result, nil
}
